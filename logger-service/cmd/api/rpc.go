package main

import (
	"context"
	"log"
	"log-service/data"
	"time"
)

type RPCServer struct{}
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	logSnippet := "\n[logger-service][rpc][LogInfo] =>"

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	
	if err != nil {
		log.Printf("%s (ERROR-client.Database.Collection): %s", logSnippet, err.Error())
		return err
	}
	log.Printf("%s (SUCCESS-client.Database.Collection)", logSnippet)

	*resp = "Proccessed payload (" + payload.Name + ") via RPC"

	return nil
}
