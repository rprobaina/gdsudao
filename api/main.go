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

type gd struct {
	data  time.Time
	gd    float64
	fonte string
	done  bool
}

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
	dataFinalISO = dataFinalISO.AddDate(0, 0, 1)

	// Consulta do banco de dados
	query := bson.M{"codigoINMET": codigoINMET, "dataMedicao": bson.M{"$gt": dataInicialISO, "$lte": dataFinalISO}} // Limitar data de inicio e data de fim

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

// calcularGrausDia tmax e tmin e calcula graus dia
func calcularGrausDia(tMin float64, tMax float64) float64 {
	TB := 10.0
	tMedia := (tMin + tMax) / 2
	gd := tMedia - TB
	return gd
}

func somaTermica(gdSlice []float64) float64 {
	somaTermica := 0.0
	for _, gd := range gdSlice {
		somaTermica += gd
	}
	return somaTermica
}

// getDiarios retorna os dados de medições diárias coletados por uma estação meteorológica
func getGrausDia(w http.ResponseWriter, r *http.Request) {

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]
	dataInicial := vars["dataInicial"]
	dataFinal := vars["dataFinal"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collectionDiarios := mongoClient.Database("gdsudao").Collection("diarios")
	collectionPrevisoes := mongoClient.Database("gdsudao").Collection("previsoes")
	collectionNormais := mongoClient.Database("gdsudao").Collection("normais")
	collectionEstacoes := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Conversao das datas para o formato ISO
	layoutISO := "2006-01-02"
	dataInicialISO, _ := time.Parse(layoutISO, dataInicial)
	dataFinalISO, _ := time.Parse(layoutISO, dataFinal)
	dataFinalISO = dataFinalISO.AddDate(0, 0, 1)

	// Cria o slice de datas
	var gds []gd
	for currentDate := dataInicialISO; currentDate != dataFinalISO; currentDate = currentDate.AddDate(0, 0, 1) {
		gd := gd{currentDate, 0, "nil", false}
		gds = append(gds, gd)
	}

	for p, x := range gds {
		var tmin, tmax float64
		var fonte string
		/*
			dataISO := x.data.Format(layoutISO)
				fmt.Println("--- --- --- --- ---")
				fmt.Println("Data: " + dataISO)
				fmt.Println("Graus Dia: " + strconv.FormatFloat(x.gd, 'E', -1, 64))
				fmt.Println("Fonte: " + x.fonte)
				fmt.Println("Done: " + strconv.FormatBool(x.done))
				fmt.Println("--- --- --- --- ---")
		*/

		// Buscar dados diários
		//max, min, err getDiario(data)
		// Consulta do banco de dados
		queryDiarios := bson.M{"codigoINMET": codigoINMET, "dataMedicao": x.data}
		var diario bson.M
		err := collectionDiarios.FindOne(context.TODO(), queryDiarios).Decode(&diario)

		if err != nil {
			// Tentar previsoes
			fmt.Fprintf(w, "Erro")
			fmt.Println("Nao achou diarios	" + x.data.Format(layoutISO))
			//findOptions := options.Find()
			queryOptions := options.FindOneOptions{}
			queryOptions.SetSort(bson.M{"dataAtualizacao": -1, "last_error_time": 1})
			//findOptions.SetSort(bson.D{{"dataAtualizacao", -1}})
			queryPrevisoes := bson.M{"codINMET": codigoINMET, "dataPrevisao": x.data}
			var previsao bson.M
			err := collectionPrevisoes.FindOne(context.TODO(), queryPrevisoes, &queryOptions).Decode(&previsao)
			if err != nil {
				// Tentar previsoes
				fmt.Fprintf(w, "Erro")
				fmt.Println("Nao achou em previsoes")
				// *** Pegando dado de NORMAIS ***
				// Codigo INMET não bate pq as normais sao com estacoes manuais
				queryEstacao := bson.M{"codigoINMET": codigoINMET}
				fmt.Println(codigoINMET)
				var estacao bson.M
				err := collectionEstacoes.FindOne(context.TODO(), queryEstacao).Decode(&estacao)
				fmt.Println(codigoINMET)
				fmt.Println(estacao)
				if err != nil {
					// Erro ao buscar estacoes
					fmt.Println("Estacao para normais nao encontarada")
					fmt.Println(err)
				} else {
					nomeEstacao := estacao["nomeEstacao"].(string)
					queryNormais := bson.M{"nomeEstacao": nomeEstacao}
					var normais bson.M
					err := collectionNormais.FindOne(context.TODO(), queryNormais).Decode(&normais)
					if err != nil {
						// Erro ao buscar normais
					} else {
						// TODO: calcular o retorno das normais. Talvezes colocar num mapa ou fazer uma função
						fmt.Println(normais)
					}
					// Buscar normais
				}

			} else {
				// *** Pegando dado de PREVISAO ***
				tmin = previsao["temperaturaMinima"].(float64)
				tmax = previsao["temperaturaMaxima"].(float64)
				fonte = "previsao"
				// Debug
				fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
					" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)
				// Atualiza o vetor de graus dia
				//fmt.Println("Previsoes")
				//fmt.Println(previsao)
			}

		} else {
			// *** Pegando dados DIÁRIOS ***
			tmin = diario["temperaturaMinima"].(float64)
			tmax = diario["temperaturaMaxima"].(float64)
			fonte = "diario"
			// Debug
			fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
				" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)

		}
		// Atualiza o vetor de graus dia
		grausDia := calcularGrausDia(tmin, tmax)
		if grausDia > 0 {
			x.gd = grausDia
			x.fonte = fonte
			x.done = true
			gds[p] = x // Atualiza o slice
		} else {
			x.done = false
		}
	}
	//fmt.Println("+++++++++++++++++++++++++++++++++++++")
	fmt.Println(gds)

	/*
		TODO: validar dados, calcular graus dia em gds, retornar valor e percentuais
	*/

	// V NADA PRESTA
	//----------------
	// Consulta do banco de dados
	queryDiarios := bson.M{"codigoINMET": codigoINMET, "dataMedicao": bson.M{"$gt": dataInicialISO, "$lte": dataFinalISO}} // Limitar data de inicio e data de fim

	// Realiza a consulta no banco de dados e retorna o valor encontrado
	var diarios []bson.M
	var gdSlice []float64
	cur, err := collectionDiarios.Find(context.Background(), queryDiarios)
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
		tMin := result["temperaturaMinima"]
		tMax := result["temperaturaMaxima"]

		gd := calcularGrausDia(tMin.(float64), tMax.(float64))

		gdSlice = append(gdSlice, gd)
	}
	somaTermica := somaTermica(gdSlice)

	if err := cur.Err(); err != nil {
		fmt.Fprintf(w, "Codigo de estação inválido")
	} else {
		json.NewEncoder(w).Encode(somaTermica)
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
	myRouter.HandleFunc("/grausdia/{codigoINMET}/{dataInicial}/{dataFinal}", getGrausDia).Methods("GET")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

//	main é função Principal do programa
func main() {
	fmt.Println("GDSudão API: on")
	defer fmt.Println("GDSudão API: off")
	handleRequests()
}
