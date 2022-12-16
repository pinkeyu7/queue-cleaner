package main

import (
	"log"
	"queue-cleaner/config"
	"queue-cleaner/queue_management"
)

func main() {
	config.InitEnv()

	_, err := queue_management.InactiveQueueList()
	if err != nil {
		log.Panic(err)
	}
}
