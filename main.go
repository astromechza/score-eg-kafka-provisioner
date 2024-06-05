package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	host, port, topicName := os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"), os.Getenv("KAFKA_TOPIC")
	if host == "" {
		log.Fatal("KAFKA_HOST is empty or not set")
	} else if port == "" {
		log.Fatal("KAFKA_PORT is empty or not set")
	} else if topicName == "" {
		log.Fatal("KAFKA_TOPIC is empty or not set")
	}

	log.Printf("connecting to %s:%s/%s and sending 'ping' every 5s", host, port, topicName)
	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%s", host, port), topicName, 0)
	if err != nil {
		log.Fatalf("failed to dial to kafka: %v", err)
	}
	defer conn.Close()

	for {
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
		if r, err := conn.Write([]byte("ping")); err != nil {
			log.Printf("error: failed to ping: %v", err)
		} else {
			slog.Info("successfully sent message", slog.Int("#bytes", r))
		}
		time.Sleep(time.Second * 5)
	}
}
