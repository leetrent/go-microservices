package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	//mongoURL = "mongodb://localhost:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

func main() {
	// Connect to mongodb
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Println(err)
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
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Start web server
	//go app.serve()

	log.Println("Starting logger-service on port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}
// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Panic(err)
// 	}
// }

func connectToMongo() (*mongo.Client, error) {

	// Create Connection Options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongodb:", err)
		return nil, err
	}

	log.Println("Connected to mongodb!")

	return c, nil
}
