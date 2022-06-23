package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ValidCollection() {
	var db = client.Database("my_db")

	// var collection = db.Collection("valid_collection")

	err := db.CreateCollection(context.Background(), "valid_collection",
		options.CreateCollection().SetValidator(bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"name", "status"},
				"properties": bson.M{
					"name": bson.M{
						"bsonType":    "string",
						"description": "must be a string and is required",
					},
					"status": bson.M{
						"enum":        []string{"INIT", "DEL"},
						"description": "can only be one of the enum values and is required",
					},
				},
			},
		}),
	)

	if err != nil {
		panic(err)
	}

	var collection = db.Collection("valid_collection")

	const (
		INIT = "INIT"
		DEL  = "DEL"
	)

	type Info struct {
		Name   string `bson:"name,omitempty"`
		Status string `bson:"status,omitempty"`
	}

	ret, err := collection.InsertOne(context.Background(), Info{
		Name:   "ganyyy",
		Status: INIT,
	})
	if err != nil {
		log.Fatalf("insert error:%v", err)
	}
	var id, _ = primitive.ObjectIDFromHex(ret.InsertedID.(primitive.ObjectID).Hex())
	log.Printf("id:%v", id)
	var info Info
	err = collection.FindOne(context.Background(), bson.M{
		"_id": id,
	}).Decode(&info)
	if err != nil {
		log.Fatalf("find one error:%v", err)
	}
	log.Printf("ret:%+v", info)
}
