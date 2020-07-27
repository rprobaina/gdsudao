// Package utils is a simple API that contains the most common database operations
package utils

import (
	"mongoapi"
	"os/exec"

	"go.mongodb.org/mongo-driver/mongo"
)

// wget is used to download a file from <url> parameter and save it to <filename>
func wget(url string, filename string) {
	exec.Command("wget", url, "-O", filename).Run()
}

// rm is used to remove the file <filename> from the filesystem
func rm(filename string) {
	exec.Command("rm", filename).Run()
}

// Open DB collection
func openConectionDB(dbName string, collectionName string) *mongo.Collection {
	dataBaseURI := "mongodb://127.0.0.1:27017"

	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database(dbName).Collection(collectionName)
	//defer mongoapi.CloseConnection(*mongoClient)

	return collection
}
