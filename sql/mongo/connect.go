package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitClient() {
	var err error
	if client, err = mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://localhost:27017").
			SetConnectTimeout(5*time.Second).SetAuth(options.Credential{
			Username: "admin",
			Password: "123456",
		})); err != nil {
		log.Fatalf("connect error:%v", err)
	}
}
