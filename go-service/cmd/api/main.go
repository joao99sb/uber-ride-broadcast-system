package main

import (
	"go-service/internal"
	"go-service/pkg/queue"
	"go-service/pkg/server"

	"github.com/joho/godotenv"
)

func Init() {
	// err := godotenv.Load("../../../.env")
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func main() {

	Init()

	queue := queue.NewQueue()
	queue.Connect()
	defer queue.Close()

	service := internal.NewService(queue)
	controller := internal.NewController(service)
	server := server.NewStreamServer(controller)

	server.Run()

}
