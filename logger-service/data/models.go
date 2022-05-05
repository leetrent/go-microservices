package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	logSnippet := "\n[logger-service][models][Insert()] =>"

	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Printf("%s (ERROR-collection.InsertOne()): %s", logSnippet, err.Error())
		return err
	}
	log.Printf("%s (SUCCESS-collection.InsertOne())", logSnippet)

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	logSnippet := "\n[logger-service][models][All] =>"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Printf("%s (ERROR-collection.Find()): %s", logSnippet, err.Error())
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Printf("%s (ERROR-cursor.Decode()): %s", logSnippet, err.Error())
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}

	log.Printf("%s (SUCCESS-len(logs)): %d", logSnippet, len(logs))

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	logSnippet := "\n[logger-service][models][GetOne] =>"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("%s (ERROR-primitive.ObjectIDFromHex()): %s", logSnippet, err.Error())
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)
	if err != nil {
		log.Printf("%s (ERROR-collection.FindOne()): %s", logSnippet, err.Error())
		return nil, err
	}
	log.Printf("%s (SUCCESS-collection.FindOne())", logSnippet)

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	logSnippet := "\n[logger-service][models][DropCollection] =>"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		log.Printf("%s (ERROR-collection.Drop()): %s", logSnippet, err.Error())
		return err
	}
	log.Printf("%s (SUCCESS-collection.Drop())", logSnippet)

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	logSnippet := "\n[logger-service][models][Update] =>"

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Printf("%s (ERROR-primitive.ObjectIDFromHex): %s", logSnippet, err.Error())
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updated_at", time.Now()},
			}},
		},
	)

	if err != nil {
		log.Printf("%s (ERROR-collection.UpdateOne): %s", logSnippet, err.Error())
		return nil, err
	}
	log.Printf("%s (SUCCESS-collection.UpdateOne)", logSnippet)

	return result, err
}
