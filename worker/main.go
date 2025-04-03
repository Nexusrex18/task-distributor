package main

import (
	"encoding/json"
	"log"

	"github.com/Nexusrex18/task-distributer/worker/nats"
	"github.com/Nexusrex18/task-distributer/worker/processor"
	"github.com/Nexusrex18/task-distributer/worker/storage"
	natsio "github.com/nats-io/nats.go"
)

type Task struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Data   []byte `json:"data"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal("NATS connection failed: ", err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("tasks", func(msg *natsio.Msg) {
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
		}
	})

	if err != nil {
		log.Fatal("Subscription failed: ", err)
	}

	select {}
}
