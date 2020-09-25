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
)

// homePage retorna as informações gerais da API
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> GDSudao API </h1>")
	fmt.Fprintf(w, "<p> Essa é uma API utilizada pelo projeto GDSusão. </p>")
	fmt.Fprintf(w, "<p> Como usar: </p>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<p> Retorna dados da estação meteorológica mais próxima: </p>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/estacao/maisproxima/{latitude}/{longitude} </ul>")
	fmt.Fprintf(w, "<p> Retorna dados de normais climatológicas de uma estação: </p>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/normais/{nomeEstacao} </ul>")
	fmt.Fprintf(w, "<p> Retorna dados diarios de uma estação: </p>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/diarios/{codigoINMET}/{dataInicial}/{dataFinal} </ul>")
	fmt.Fprintf(w, "<p> Retorna dados de previsão do tempo de uma estação: </p>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/previsoes/{codigoINMET}/{dataAtual} </ul>")
}

// getNearStation retorna os dados da estação meteorológica mais proxima
func getNearStation(w http.ResponseWriter, r *http.Request) {

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	latitude := vars["latitude"]
	longitude := vars["longitude"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Converte as coordenadas para float64
	lat, _ := strconv.ParseFloat(latitude, 64)
	lon, _ := strconv.ParseFloat(longitude, 64)

	// Consulta geoespacial
	query := bson.D{
		{"localizacao.coordenadas", bson.D{
			{"$near", bson.D{
				{"$geometry", bson.D{
					{"type", "Point"},
					{"coordinates", bson.A{lat, lon}}},
				}},
			}},
		}}

	// Realiza a consulta no banco de dados e retorna o valor encontrado
	var resultado bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&resultado)

	if err != nil {
		fmt.Fprintf(w, "Coordenadas inválidas")
	} else {
		json.NewEncoder(w).Encode(resultado)
	}
}

// getNormais retorna os dados das normais climatológicas de uma estação meteorológica
func getNormais(w http.ResponseWriter, r *http.Request) {

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoEstacao := vars["nomeEstacao"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("normais")
	defer mongoapi.CloseConnection(*mongoClient)

	// Consulta do banco de dados
	query := bson.M{"nomeEstacao": codigoEstacao}

	// Realiza a consulta no banco de dados e retorna o valor encontrado
	var normais bson.M
	err := collection.FindOne(context.TODO(), query).Decode(&normais)

	if err != nil {
		fmt.Fprintf(w, "Estação inexistente")
	} else {
		json.NewEncoder(w).Encode(normais)
	}
}

// getDiarios retorna os dados de medições diárias coletados por uma estação meteorológica
func getDiarios(w http.ResponseWriter, r *http.Request) {

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]
	dataInicial := vars["dataInicial"]
	dataFinal := vars["dataFinal"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("diarios")
	defer mongoapi.CloseConnection(*mongoClient)

	// Conversao das datas para o formato ISO
	layoutISO := "2006-01-02"
	dataInicialISO, _ := time.Parse(layoutISO, dataInicial)
	dataFinalISO, _ := time.Parse(layoutISO, dataFinal)

	// Consulta do banco de dados
	query := bson.M{"codigoINMET": codigoINMET, "dataMedicao": bson.M{"$gte": dataInicialISO, "$lte": dataFinalISO}} // Limitar data de inicio e data de fim

	// Realiza a consulta no banco de dados e retorna o valor encontrado
	var diarios []bson.M
	cur, err := collection.Find(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		diarios = append(diarios, result)
	}
	if err := cur.Err(); err != nil {
		fmt.Fprintf(w, "Codigo de estação inválido")
	} else {
		json.NewEncoder(w).Encode(diarios)
	}
}

// getPrevisoes retorna os dados de previsão de tempo gerados para uma localidade
func getPrevisoes(w http.ResponseWriter, r *http.Request) {

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]
	dataAtual := vars["dataAtual"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("previsoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Conversao das datas para o formato ISO
	layoutISO := "2006-01-02"
	data, _ := time.Parse(layoutISO, dataAtual)

	// Consulta do banco de dados
	query := bson.M{"codINMET": codigoINMET, "dataAtualizacao": data}

	var previsoes []bson.M
	cur, err := collection.Find(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	// Realiza a consulta no banco de dados e retorna o valor encontrado
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		var previsao bson.M
		err := cur.Decode(&previsao)
		if err != nil {
			//log.Fatal(err)
		}
		previsoes = append(previsoes, previsao)
	}
	if err := cur.Err(); err != nil {
		fmt.Fprintf(w, "Codigo de localidade inválido")
	} else {
		json.NewEncoder(w).Encode(previsoes)
	}

}

//	handleRequests trata das requisições (mapeia a requisição para a função adequada)
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/estacao/maisproxima/{latitude}/{longitude}", getNearStation).Methods("GET")
	myRouter.HandleFunc("/normais/{nomeEstacao}", getNormais).Methods("GET")
	myRouter.HandleFunc("/diarios/{codigoINMET}/{dataInicial}/{dataFinal}", getDiarios).Methods("GET")
	myRouter.HandleFunc("/previsoes/{codigoINMET}/{dataAtual}", getPrevisoes).Methods("GET")

	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

//	main é função Principal do programa
func main() {
	fmt.Println("GDSudão API: on")
	defer fmt.Println("GDSudão API: off")
	handleRequests()
}
