package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongoapi"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// homePage returns a simple message of this API
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a simple REST API writen in Golang that uses MongoDB as database.")
}

// getNearStation retorna os dados da estação meteorológica mais proxima
func getNearStation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	latitude := vars["latitude"]
	longitude := vars["longitude"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	lat, _ := strconv.ParseFloat(latitude, 64)
	lon, _ := strconv.ParseFloat(longitude, 64)

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"
	query := bson.D{
		{"localizacao.coordenadas", bson.D{
			{"$near", bson.D{
				{"$geometry", bson.D{
					{"type", "Point"},
					{"coordinates", bson.A{lat, lon}}},
				}},
			}},
		}}

	var resultado bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&resultado)

	if err != nil {
		fmt.Println("Localização invalida")
	} else {
		json.NewEncoder(w).Encode(resultado)
	}
}

// getNormais retorna os dados da estação meteorológica mais proxima
func getNormais(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	codigoEstacao := vars["nomeEstacao"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("normais")
	defer mongoapi.CloseConnection(*mongoClient)

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"
	//query := bson.M{"codigoINMET": codigoEstacao} Nao é o mesmo codigo de uma estacao automatica

	query := bson.M{"nomeEstacao": codigoEstacao}

	var normais bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&normais)

	fmt.Println(err)
	fmt.Println(normais)

	if err != nil {
		fmt.Println("Localização invalida")
	} else {
		json.NewEncoder(w).Encode(normais)
	}
}

// getDiarios
func getDiarios(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("diarios")
	defer mongoapi.CloseConnection(*mongoClient)

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"
	//query := bson.M{"codigoINMET": codigoEstacao} Nao é o mesmo codigo de uma estacao automatica

	query := bson.M{"codigoINMET": codigoINMET}

	cur, err := collection.Find(context.Background(), query)
	if err != nil {
		//log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result...
		fmt.Println(result)
		// To get the raw bson bytes use cursor.Current
		//raw := cur.Current
		// do something with raw...
	}
	if err := cur.Err(); err != nil {
		fmt.Println("Localização invalida")
	} else {

		//json.NewEncoder(w).Encode(result)
	}
}

// getPrevisoes
func getPrevisoes(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	codigoEstacao := vars["codigoINMET"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("previsoes")
	defer mongoapi.CloseConnection(*mongoClient)

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"
	//query := bson.M{"codigoINMET": codigoEstacao} Nao é o mesmo codigo de uma estacao automatica

	query := bson.M{"nomeEstacao": codigoEstacao}

	var normais bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&normais)

	fmt.Println(err)
	fmt.Println(normais)

	if err != nil {
		fmt.Println("Localização invalida")
	} else {
		json.NewEncoder(w).Encode(normais)
	}
}

//	Trata das requisições (mapeia a requisição para a função adequada)
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/estacao/maisproxima/{latitude}/{longitude}", getNearStation).Methods("GET")
	myRouter.HandleFunc("/normais/{nomeEstacao}", getNormais).Methods("GET")
	myRouter.HandleFunc("/diarios/{codigoINMET}/{data_inicio}/{data_fim}", getDiarios).Methods("GET")
	myRouter.HandleFunc("/previsoes/{codigoCPTEC}/{data_inicio}/{data_fim}", getPrevisoes).Methods("GET")

	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

//	Função Principal do programa
func main() {
	fmt.Println("API: on")
	defer fmt.Println("API: off")
	handleRequests()
}
