package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nexusrex18/task-distributer/api/handlers"
	"github.com/Nexusrex18/task-distributer/api/metrics"
	"github.com/Nexusrex18/task-distributer/api/nats"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// gin.SetMode(gin.ReleaseMode) //ths is for production

	r := gin.New()
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[API] %s | %3d | %13v | %15s | %-7s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	}), gin.Recovery(),
	)

	_ = r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

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

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.POST("/tasks", handlers.HandleTaskSubmission(nc))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		nc.Close()
		os.Exit(0)
	}()

	log.Fatal(r.Run(":8080"))
}
