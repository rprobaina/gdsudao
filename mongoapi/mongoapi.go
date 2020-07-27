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

// User is data type of a simple user used for this example.
type User struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

// StartConnection starts a database connection.
func StartConnection(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Println("Error starting the database connection: ", err)
	} else {
		fmt.Println("Database connection successfully started!")
	}

	return client
}

// CloseConnection close the the database connection.
func CloseConnection(client mongo.Client) error {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Println("Error closing the database connection: ", err)
	} else {
		fmt.Println("Connection closed!")
	}
	return err

}

// InsertUser insert a user in a collection.
func InsertUser(client mongo.Client, collection mongo.Collection, user User) error {
	result, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		log.Println("404 - ", err)
	} else {
		fmt.Println("User successfully inserted! ", result.InsertedID)
	}

	return err
}

// QueryUser returns a user matches with the query.
func QueryUser(client mongo.Client, collection mongo.Collection, query bson.D) (User, error) {
	var result User
	err := collection.FindOne(context.TODO(), query).Decode(&result)

	if err != nil {
		log.Println("404 - User not found", result)
		result = User{"", "", 0}
	} else {
		fmt.Println("User successfully consulted! ", result)
	}

	return result, err
}

// QueryUsers returns a slice of users that matches with the query.
func QueryUsers(client mongo.Client, collection mongo.Collection, query bson.D) ([]User, error) {
	cursor, err := collection.Find(context.TODO(), query)

	if err != nil {
		log.Println("400 - ", err)
	} else {
		fmt.Println("Users successfully consulted!")
	}

	var users []User

	for cursor.Next(context.TODO()) {
		var user User

		// decode the document
		err := cursor.Decode(&user)

		if err != nil {
			log.Println("400 - ", err)
		}

		users = append(users, user)
	}

	return users, err
}

// UpdateUser updates a user in a collection.
func UpdateUser(client mongo.Client, collection mongo.Collection, query bson.D, update bson.D) error {
	result, err := collection.UpdateMany(context.TODO(), query, update)

	if err != nil {
		log.Println("400 - ", err)
	} else {
		fmt.Println("User updated!", result)
	}

	return err
}

// DeleteUser delete a user from a collection.
func DeleteUser(client mongo.Client, collection mongo.Collection, query bson.D) error {
	result, err := collection.DeleteOne(context.TODO(), query)

	if err != nil {
		log.Println("400 - ", err)
	} else {
		fmt.Println("User deleted!", result)
	}

	return err
}

// ClearCollection erase all documents from a collection.
func ClearCollection(client mongo.Client, collection mongo.Collection) error {
	query := bson.D{{}}
	result, err := collection.DeleteMany(context.TODO(), query)

	if err != nil {
		log.Println("400", err)
	} else {
		fmt.Println("All documents are gone!", result)
	}

	return err
}
