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

type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	logSnippet := "\n[logger-service][grpc][WriteLog] =>"

	input := req.GetLogEntry()

	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		log.Printf("%s (ERROR - l.Models.LogEntry.Insert): %s", logSnippet, err.Error())
		res := &logs.LogResponse{
			Result: "Log entry insert failed",
		}
		return res, err
	}
	log.Printf("%s (SUCCESS - l.Models.LogEntry.Insert)", logSnippet)

	res := &logs.LogResponse{
		Result: "Log entry insert succeeded",
	}

	return res, nil
}

func (app *Config) gRPCListen() {
	logSnippet := "\n[logger-service][grpc][gRPCListen] =>"

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		errMsg := fmt.Sprintf("%s (ERROR - net.Listen): %s", logSnippet, err.Error())
		log.Println(errMsg)
		log.Fatalf(errMsg)
	}
	log.Printf("%s (SUCCESS - net.Listen)", logSnippet)

	srv := grpc.NewServer()
	logs.RegisterLogServiceServer(srv, &LogServer{Models: app.Models})

	log.Printf("%s (INFO - gRPC Server started on port %s)", logSnippet, gRpcPort)

	if err := srv.Serve(lis); err != nil {
		errMsg := fmt.Sprintf("%s (ERROR - srv.Serve): %s", logSnippet, err.Error())
		log.Println(errMsg)
		log.Fatalf(errMsg)
	}
	log.Printf("%s (SUCCESS - srv.Serve)", logSnippet)
}
