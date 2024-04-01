package db

import (
	"log"
	"os"

	"example.com/m/social_media/utils"
	"github.com/hibiken/asynq"
)

var REDISCLIENT *asynq.Client
var SRV *asynq.Server

// InitREDIS initializes the redis client
func InitREDIS() {
	REDIS_URL := os.Getenv("REDIS_URL")

	REDISCLIENT = asynq.NewClient(asynq.RedisClientOpt{Addr: REDIS_URL})
}

// HandleServer sets up the server to listen for tasks
func HandleServer() {
	REDIS_URL := os.Getenv("REDIS_URL")
	SRV = asynq.NewServer(
		asynq.RedisClientOpt{Addr: REDIS_URL},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default": 3,
				"low": 1,
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(utils.TypeEmailDelivery, utils.HandleEmailDeliveryTask)
	if err := SRV.Run(mux); err != nil {
		log.Printf("Could not run server: %v", err)
	}
}