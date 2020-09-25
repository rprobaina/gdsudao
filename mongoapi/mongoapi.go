// Package mongoapi is a simple API that contains the most common database operations
package mongoapi

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StartConnection starts a database connection.
func StartConnection(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error starting the database connection: ", err)
	}
	return client
}

// CloseConnection close the the database connection.
func CloseConnection(client mongo.Client) error {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Println("Error closing the database connection: ", err)
	}
	return err

}

// InsertUser insert a user in a collection.
func InsertDocument(client mongo.Client, collection mongo.Collection, document bson.D) error {
	result, err := collection.InsertOne(context.TODO(), document)

	if err != nil {
		log.Println("404 - ", err)
	} else {
		fmt.Println("Document successfully inserted! ", result.InsertedID)
	}
	return err
}
