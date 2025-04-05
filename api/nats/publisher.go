package nats

import (
	"time"
	"log"

	"github.com/nats-io/nats.go"
)

func Connect(url string) (*nats.Conn, error) {
	nc, err := nats.Connect(url,
		nats.Name("Task API"),
		nats.PingInterval(20*time.Second),
		nats.MaxReconnects(5),
	)
	if err == nil {
		log.Println("Connected to NATS at:",url)
	}

	return nc,err
}
