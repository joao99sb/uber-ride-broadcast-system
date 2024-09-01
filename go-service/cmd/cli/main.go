package main

import (
	"go-service/pkg/queue"
	"log"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load("./../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}
func main() {
	Init()

	q := queue.NewQueue()
	q.Connect()
	defer q.Close()

	defaultDestination := q.BuildDefaultDestinationQueue()
	err := q.CreateQueueAndBind(defaultDestination)
	if err != nil {
		log.Fatalf("Error creating and binding queue: %v", err)
	}
}
