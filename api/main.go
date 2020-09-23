package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongoapi"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Station tipo de dado que contem os codigos do CPTEC e do INMET de um estação
type Station struct {
	DataAtualizacao string
}

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
	dataInicial := vars["dataInicial"]
	dataFinal := vars["dataFinal"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("diarios")
	defer mongoapi.CloseConnection(*mongoClient)

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"
	//query := bson.M{"codigoINMET": codigoEstacao} Nao é o mesmo codigo de uma estacao automatica

	// Consulta testada no compas
	// {codigoINMET: "A827", dataMedicao: {"$gte": ISODate('2020-09-17T00:00:00.000Z'), "$lte": ISODate('2020-09-20T00:00:00.000Z')}}

	// Conversao das datas para o formato ISO
	layoutISO := "2006-01-02"
	dataInicialISO, _ := time.Parse(layoutISO, dataInicial)
	dataFinalISO, _ := time.Parse(layoutISO, dataFinal)

	query := bson.M{"codigoINMET": codigoINMET, "dataMedicao": bson.M{"$gte": dataInicialISO, "$lte": dataFinalISO}} // Limitar data de inicio e data de fim

	var diarios []bson.M
	cur, err := collection.Find(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result...
		fmt.Println(result)
		diarios = append(diarios, result)

		// To get the raw bson bytes use cursor.Current
		//raw := cur.Current
		// do something with raw...
	}
	if err := cur.Err(); err != nil {
		fmt.Println("Localização invalida")
	} else {
		json.NewEncoder(w).Encode(diarios)
	}
}

// getPrevisoes
func getPrevisoes(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	codigoCPTEC := vars["codigoCPTEC"]
	dataInicial := vars["dataInicial"]
	dataFinal := vars["dataFinal"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("previsoes")
	defer mongoapi.CloseConnection(*mongoClient)

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"
	//query := bson.M{"codigoINMET": codigoEstacao} Nao é o mesmo codigo de uma estacao automatica

	// {dataPrevisao: {"$gte": ISODate('2020-09-17T00:00:00.000Z')}}

	//query := bson.M{"codCPTEC": codigoCPTEC} // Limitar data de inicio e data de fim

	// Conversao das datas para o formato ISO
	layoutISO := "2006-01-02"
	dataInicialISO, _ := time.Parse(layoutISO, dataInicial)
	dataFinalISO, _ := time.Parse(layoutISO, dataFinal)

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"dataAtualizacao": -1})
	query := bson.M{"codCPTEC": codigoCPTEC, "dataPrevisao": bson.M{"$gte": dataInicialISO, "$lte": dataFinalISO}} // Limitar data de inicio e data de fim

	var previsoes []bson.M
	cur, err := collection.Find(context.Background(), query, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	ultimaData := "0000"
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var previsao bson.M
		err := cur.Decode(&previsao)
		if err != nil {
			//log.Fatal(err)
		}

		p := previsao["dataAtualizacao"]
		dataAtual := fmt.Sprintf("%v", p)

		// Não esta correto!
		if ultimaData != dataAtual {
			ultimaData = dataAtual
			fmt.Println(previsao)
			previsoes = append(previsoes, previsao)
		}

		// To get the raw bson bytes use cursor.Current
		//raw := cur.Current
		// do something with raw...
	}
	if err := cur.Err(); err != nil {
		fmt.Println("Localização invalida")
	} else {
		json.NewEncoder(w).Encode(previsoes)
	}

}

//	Trata das requisições (mapeia a requisição para a função adequada)
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/estacao/maisproxima/{latitude}/{longitude}", getNearStation).Methods("GET")
	myRouter.HandleFunc("/normais/{nomeEstacao}", getNormais).Methods("GET")
	myRouter.HandleFunc("/diarios/{codigoINMET}/{dataInicial}/{dataFinal}", getDiarios).Methods("GET")
	myRouter.HandleFunc("/previsoes/{codigoCPTEC}/{dataInicial}/{dataFinal}", getPrevisoes).Methods("GET")

	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

//	Função Principal do programa
func main() {
	fmt.Println("API: on")
	defer fmt.Println("API: off")
	handleRequests()
}
