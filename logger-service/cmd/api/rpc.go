package main

import (
	"context"
	"fmt"
	"log-service/data"
	"time"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are avaiable
// over RPC, as long as they are exported.
type RPCServer struct {
}

// the type of data that we're going to receive from RPC for RPCServer
type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		fmt.Println("error writing to mongo:", err)
		return err
	}
	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}
