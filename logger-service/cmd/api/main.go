package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	//mongoURL = "mongodb://localhost:27017"
	//gRpcPort = "50001"
)

var client *mongo.Client

func main() {
	logSnippet := "\n[logger-service][main][main] =>"
	log.Printf("%s (ENTRY-POINT):", logSnippet)

	// Connect to mongodb
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Printf("%s (ERROR-connectToMongo): %s", logSnippet, err.Error())
		log.Panic(err)
	}
	client = mongoClient

	// Create a context in order to disconnect from mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Close mongodb connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println("Error encountered when attempting to disconnect from mongodb")
			log.Printf("%s (ERROR-client.Disconnect): %s", logSnippet, err.Error())
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register RPC Server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	// Start web server
	log.Printf("%s (Starting logger-service on port %s)", logSnippet, webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("%s (ERROR-srv.ListenAndServe): %s", logSnippet, err.Error())
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	logSnippet := "\n[logger-service][main][rpcListen()] =>"
	log.Println("Starting RPC server on port ", rpcPort)

	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		log.Printf("%s (ERROR-netListen): %s", logSnippet, err.Error())
		return err
	}
	defer listen.Close()
	log.Printf("%s (SUCCESS-netListen):", logSnippet)

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			log.Printf("%s (CONTINUE-listen.Accept): %s", logSnippet, err.Error())
			continue
		}
		log.Printf("%s (SUCCESS-listen.Accept)", logSnippet)

		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	logSnippet := "\n[logger-service][main][connectToMongo()] =>"

	// Create Connection Options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Printf("%s (ERROR-mongo.Connect()): %s", logSnippet, err.Error())
		return nil, err
	}

	log.Printf("%s (SUCCESS-mongo.Connect()): ", logSnippet)

	return c, nil
}
