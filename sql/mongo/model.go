package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimePrint struct {
	StartTime int64 `bson:"start_time,omitempty"`
	EndTime   int64 `bson:"end_time,omitempty"`
}

type LogRecord struct {
	JobName string    `bson:"job_name,omitempty"`
	Command string    `bson:"command,omitempty"`
	Err     string    `bson:"err,omitempty"`
	Content string    `bson:"content,omitempty"`
	Tp      TimePrint `bson:"tp,omitempty"`
}

func InsertOne() {
	var (
		db         *mongo.Database
		collection *mongo.Collection
	)

	db = client.Database("my_db")

	collection = db.Collection("my_collection")

	var id string
	if result, err := collection.InsertOne(
		context.Background(), LogRecord{
			JobName: "123",
			Command: "456",
			Err:     "123123",
			Content: "1232131",
			Tp: TimePrint{
				StartTime: 0,
				EndTime:   10000,
			},
		}); err != nil {
		log.Fatalf("insert error:%v", err)
	} else {
		id = result.InsertedID.(primitive.ObjectID).Hex()
		log.Printf("insert success! id:%s", string(id[:]))
	}

	var hexId, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatalf("prase objectId error:%v", err)
	}

	var logRecord LogRecord
	cursor := collection.FindOne(context.Background(), bson.M{
		"_id": hexId,
	})
	if err := cursor.Err(); err != nil {
		log.Fatalf("find one error:%+v", err)
	}

	if err := cursor.Decode(&logRecord); err != nil {
		log.Fatalf("decode error:%+v", err)
	}

	log.Printf("%+v", logRecord)
}
