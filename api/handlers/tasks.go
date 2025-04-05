package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type TaskRequest struct {
	Image  []byte `form:"image" binding:"required"`
	Width  int    `form:"width" binding:"required"`
	Height int    `form:"height" binding:"required"`
}

func HandleTaskSubmission(nc *nats.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "file upload error",
				"detail": err.Error(),
			})
			return
		}

		log.Printf("Receiving file: %s (Size: %d)", file.Filename,file.Size)

		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to read image",
				"detail": err.Error(),
			})
			return
		}
		defer f.Close()

		imageData, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process image"})
			return
		}

		widhtStr := c.PostForm("width")
		heightStr := c.PostForm("height")
		var width , height int
		if _,err := fmt.Sscanf(widhtStr, "%d",&width); err!=nil || width <= 0 {
			c.JSON(http.StatusBadRequest,gin.H{"error": "invalid or missing width"})
			return
		}
		if _,err := fmt.Sscanf(heightStr, "%d", &height);err!=nil || height <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid or missing height"})
			return
		}

		task := struct {
			ID     string `json:"id"`
			Type   string `json:"type"`
			Data   []byte `json:"data"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		}{
			ID:     uuid.New().String(),
			Type:   "resize",
			Data:   imageData,
			Width:  width,
			Height: height,
		}

		taskJSON, _ := json.Marshal(task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed to encode task"})
			return
		}
		if err := nc.Publish("tasks", taskJSON); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue task"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"id":     task.ID,
			"status": "queued",
		})
	}
}
