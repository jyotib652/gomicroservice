package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"log-service/logs"
	"net"

	"google.golang.org/grpc"
)

// This is going to be the gRPC Server

type LogServer struct {
	logs.UnimplementedLogServiceServer             // This member(type) of the struct is required for every service of gRPC
	Models                             data.Models // to get the access to necessary methods to write to MongoDB
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	// And we write the log to mongodb in logger-service
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}
	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// Return response
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}

// Now, listen for gRPC connections
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC connections :%v", err)
	}

	// Now, listen for connections
	s := grpc.NewServer()

	// Now, register the service
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC Server has started on port %s\n", gRpcPort)

	// Now, do the listening
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC connections :%v", err)
	}
}
