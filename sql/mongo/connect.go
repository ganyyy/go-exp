package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitClient() {
	var err error
	if client, err = mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://localhost:27017"),
	); err != nil {
		log.Fatalf("connect error:%v", err)
	}
}
