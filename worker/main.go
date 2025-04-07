package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/Nexusrex18/task-distributer/worker/metrics"
	"github.com/Nexusrex18/task-distributer/worker/nats"
	"github.com/Nexusrex18/task-distributer/worker/processor"
	"github.com/Nexusrex18/task-distributer/worker/storage"
	natsio "github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Task struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Data   []byte `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func main() {
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://nats:4222"
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("NATS connection failed: ", err)
	}
	defer nc.Close()
	
	metrics.Register()

	go func ()  {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":9091",nil))	
	}()

	_, err = nc.Subscribe("tasks", func(msg *natsio.Msg) {
		start := time.Now()
		var task Task 
		err := json.Unmarshal(msg.Data, &task)
		if err != nil {
			log.Printf("Failed to decode task: %v", err)
			return
		}

		resized, err := processor.ResizeImage(task.Data, task.Width, task.Height)

		if err != nil {
			log.Printf("Resize failed: %v", err)
			return
		}

		err = storage.SaveToMinIO("processed", task.ID+".jpg", resized)
		if err != nil {
			log.Printf("MinIO save failed: %v", err)
		} else {
			log.Printf("Processed task %s", task.ID)
			metrics.TasksProcessed.Inc()
			metrics.ProcessingDuration.Observe(time.Since(start).Seconds())
		}
	})

	if err != nil {
		log.Fatal("Subscription failed: ", err)
	}

	select {}
}
