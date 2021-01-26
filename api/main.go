package main

import (
	"context"
	"encoding/json"
	"errors"
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
	fmt.Fprintf(w, "<h1> GD Sudão Application Programing Interface </h1>")
	fmt.Fprintf(w, "<p> Essa API compõem o FMIS GD Sudão, densenvolvido em parceiria pela Universidade Federal do Pampa e Embrapa Pecuária Sul. </p>")
	fmt.Fprintf(w, "<h2> Como usar: </h2>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<ul>")
	fmt.Fprintf(w, "<li> <b> Retorna dados da estação meteorológica mais próxima a um ponto: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/estacao/maisproxima/{latitude}/{longitude} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna dados de normais climatológicas de uma estação: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/normais/{nomeEstacao} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna dados diários de uma estação: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/diarios/{codigoINMET}/{dataInicial}/{dataFinal} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna dados de previsão do tempo de uma estação: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/previsoes/{codigoINMET}/{dataAtual} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna a soma térmica para o Capim-Sudão BRS-Stribo e a proproção de dados utililizados: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/gdsudao/{codigoINMET}/{dataInicial}/{dataFinal} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna a soma térmica para determinada temperatura basal e a proproção de dados utililizados: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/somatermica/{temperaturaBasal}/{codigoINMET}/{dataInicial}/{dataFinal} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna a data do próximo pastejo do Capim-Sudão BRS-Stribo: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080//gdsudaoProximoPastejo/{codigoINMET}/{dataInicial}/{numeroPastejos} </ul>")
	fmt.Fprintf(w, "<p> </p>")
	fmt.Fprintf(w, "<li> <b> Retorna o número estimados de pastejos de uma região: </b> </li>")
	fmt.Fprintf(w, "<ul> https://localhost:8080/gdsudao/{codigoINMET}/{dataInicial}/{dataFinal} </ul>")
	fmt.Fprintf(w, "</ul>")
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
					{"coordinates", bson.A{lon, lat}}},
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
func calcularGrausDia(tMin float64, tMax float64, TB float64) float64 {
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

func getNormal(data string, normais bson.M) float64 {
	var normal float64

	switch data {
	case "January":
		normal = normais["normalJaneiro"].(float64)
	case "February":
		normal = normais["normalFevereiro"].(float64)
	case "March":
		normal = normais["normalMarco"].(float64)
	case "April":
		normal = normais["normalAbril"].(float64)
	case "May":
		normal = normais["normalMaio"].(float64)
	case "June":
		normal = normais["normalJunho"].(float64)
	case "July":
		normal = normais["normalJulho"].(float64)
	case "August":
		normal = normais["normalAgosto"].(float64)
	case "September":
		normal = normais["normalSetembro"].(float64)
	case "October":
		normal = normais["normalOutubro"].(float64)
	case "November":
		normal = normais["normalNovembro"].(float64)
	case "December":
		normal = normais["normalDezembro"].(float64)
	default:
		normal = 0.0
	}

	return normal
}

// getDiarios retorna os dados de medições diárias coletados por uma estação meteorológica
func getGrausDiaSudao(w http.ResponseWriter, r *http.Request) {

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
	dataHoje := time.Now()
	dataDiarios := dataHoje.AddDate(0, 0, -2)
	dataPrevisoes := dataHoje.AddDate(0, 0, 14)

	// Cria o slice de datas
	var gds []gd
	for currentDate := dataInicialISO; currentDate != dataFinalISO; currentDate = currentDate.AddDate(0, 0, 1) {
		gd := gd{currentDate, 0, "nil", false}
		gds = append(gds, gd)
	}

	// Tentar previsoes
	//fmt.Fprintf(w, "Erro")
	//fmt.Println("Nao achou em previsoes")
	// *** Pegando dado de NORMAIS ***
	// Codigo INMET não bate pq as normais sao com estacoes manuais
	var normais bson.M
	queryEstacao := bson.M{"codigoINMET": codigoINMET}
	//fmt.Println(codigoINMET)
	var estacao bson.M
	err := collectionEstacoes.FindOne(context.TODO(), queryEstacao).Decode(&estacao)
	//fmt.Println(codigoINMET)
	//fmt.Println(estacao)
	if err != nil {
		// Erro ao buscar estacoes
		//fmt.Println("Estacao para normais nao encontarada")
		fmt.Println(err)
	} else {
		nomeEstacao := estacao["nomeEstacao"].(string)
		queryNormais := bson.M{"nomeEstacao": nomeEstacao}

		err := collectionNormais.FindOne(context.TODO(), queryNormais).Decode(&normais)
		if err != nil {
			normais = nil
		}
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
		var diario bson.M
		var err error = nil
		//Otimizacao
		if x.data.After(dataDiarios) {
			err = errors.New("otimizacao")
		} else {
			//fmt.Println("deveria pular")
			//fmt.Println(x.data.Format(layoutISO))
			//fmt.Println(dataHoje.Format(layoutISO))
			queryDiarios := bson.M{"codigoINMET": codigoINMET, "dataMedicao": x.data}
			err = collectionDiarios.FindOne(context.TODO(), queryDiarios).Decode(&diario)
		}

		if err != nil {
			// Tentar previsoes
			//fmt.Fprintf(w, "Erro")
			//fmt.Println("Nao achou diarios	" + x.data.Format(layoutISO) + x.data.Month().String())
			//findOptions := options.Find()

			err = nil
			var previsao bson.M
			//Otimizacao
			if x.data.After(dataPrevisoes) {
				err = errors.New("otimizacao")
			} else {
				//fmt.Println("deveria pular")
				//fmt.Println(x.data.Format(layoutISO))
				//fmt.Println(dataHoje.Format(layoutISO))
				queryOptions := options.FindOneOptions{}
				queryOptions.SetSort(bson.M{"dataAtualizacao": -1, "last_error_time": 1})
				//findOptions.SetSort(bson.D{{"dataAtualizacao", -1}})
				queryPrevisoes := bson.M{"codINMET": codigoINMET, "dataPrevisao": x.data}
				err = collectionPrevisoes.FindOne(context.TODO(), queryPrevisoes, &queryOptions).Decode(&previsao)
			}
			if err != nil {

				if normais == nil {
					// Erro ao buscar normais
				} else {
					// TODO: calcular o retorno das normais. Talvezes colocar num mapa ou fazer uma função
					//fmt.Println(normais)
					//data :=
					normal := getNormal(x.data.Month().String(), normais)
					//fmt.Println(normal)
					tmin = normal
					tmax = normal
					fonte = "normal"
				}
				// Buscar normais

			} else {
				// *** Pegando dado de PREVISAO ***
				tmin = previsao["temperaturaMinima"].(float64)
				tmax = previsao["temperaturaMaxima"].(float64)
				fonte = "previsao"
				// Debug
				//fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
				//	" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)
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
			//fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
			//	" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)

		}
		// Atualiza o vetor de graus dia
		grausDia := calcularGrausDia(tmin, tmax, 10.0)
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
	//fmt.Println(gds)

	var st = 0.0
	var qDia = 0.0
	var qPre = 0.0
	var qNor = 0.0
	var qTot = 0.0
	fmt.Println("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
	for i, g := range gds {
		fmt.Printf("Item: %d \t Data: %s \t Graus-dia: %f \t Fonte: %s \t Done: %t\n", i, g.data.Format(layoutISO), g.gd, g.fonte, g.done)
		if g.done {
			st += g.gd
			switch g.fonte {
			case "diario":
				qDia++
			case "previsao":
				qPre++
			case "normal":
				qNor++
			}
			qTot++
		} else {
			fmt.Println("Algum dado nao foi encontrado")
		}

	}
	fmt.Println("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
	fmt.Printf("Intervalo: %s - %s \t Soma Termica: %f (graus dia) \t Diários: %f%% \t Previsões: %f%% \t Normais: %f%% \n", dataInicial, dataFinal, st, ((qDia / qTot) * 100), ((qPre / qTot) * 100), ((qNor / qTot) * 100))

	/*
		TODO: validar dados, calcular graus dia em gds, retornar valor e percentuais
	*/
	pDiario := ((qDia / qTot) * 100)
	pPrevisaoes := ((qPre / qTot) * 100)
	pNormais := ((qNor / qTot) * 100)

	resposta := bson.M{"Soma Termica": st, "Diarios": pDiario, "Previsoes": pPrevisaoes, "Normais": pNormais}

	if gds == nil {
		fmt.Fprintf(w, "Codigo de estação inválido")
	} else {
		json.NewEncoder(w).Encode(resposta)
	}

}

//888
// getDiarios retorna os dados de medições diárias coletados por uma estação meteorológica
func getGrausDia(w http.ResponseWriter, r *http.Request) {

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]
	dataInicial := vars["dataInicial"]
	dataFinal := vars["dataFinal"]
	temperaturaBasal := vars["temperaturaBasal"]

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
	dataHoje := time.Now()
	dataDiarios := dataHoje.AddDate(0, 0, -2)
	dataPrevisoes := dataHoje.AddDate(0, 0, 14)

	//Conversão de string para float
	TB, _ := strconv.ParseFloat(temperaturaBasal, 64)

	// Cria o slice de datas
	var gds []gd
	for currentDate := dataInicialISO; currentDate != dataFinalISO; currentDate = currentDate.AddDate(0, 0, 1) {
		gd := gd{currentDate, 0, "nil", false}
		gds = append(gds, gd)
	}

	// Tentar previsoes
	//fmt.Fprintf(w, "Erro")
	//fmt.Println("Nao achou em previsoes")
	// *** Pegando dado de NORMAIS ***
	// Codigo INMET não bate pq as normais sao com estacoes manuais
	var normais bson.M
	queryEstacao := bson.M{"codigoINMET": codigoINMET}
	//fmt.Println(codigoINMET)
	var estacao bson.M
	err := collectionEstacoes.FindOne(context.TODO(), queryEstacao).Decode(&estacao)
	//fmt.Println(codigoINMET)
	//fmt.Println(estacao)
	if err != nil {
		// Erro ao buscar estacoes
		//fmt.Println("Estacao para normais nao encontarada")
		fmt.Println(err)
	} else {
		nomeEstacao := estacao["nomeEstacao"].(string)
		queryNormais := bson.M{"nomeEstacao": nomeEstacao}

		err := collectionNormais.FindOne(context.TODO(), queryNormais).Decode(&normais)
		if err != nil {
			normais = nil
		}
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
		var diario bson.M
		var err error = nil
		//Otimizacao
		if x.data.After(dataDiarios) {
			err = errors.New("otimizacao")
		} else {
			//fmt.Println("deveria pular")
			//fmt.Println(x.data.Format(layoutISO))
			//fmt.Println(dataHoje.Format(layoutISO))
			queryDiarios := bson.M{"codigoINMET": codigoINMET, "dataMedicao": x.data}
			err = collectionDiarios.FindOne(context.TODO(), queryDiarios).Decode(&diario)
		}

		if err != nil {
			// Tentar previsoes
			//fmt.Fprintf(w, "Erro")
			//fmt.Println("Nao achou diarios	" + x.data.Format(layoutISO) + x.data.Month().String())
			//findOptions := options.Find()

			err = nil
			var previsao bson.M
			//Otimizacao
			if x.data.After(dataPrevisoes) {
				err = errors.New("otimizacao")
			} else {
				//fmt.Println("deveria pular")
				//fmt.Println(x.data.Format(layoutISO))
				//fmt.Println(dataHoje.Format(layoutISO))
				queryOptions := options.FindOneOptions{}
				queryOptions.SetSort(bson.M{"dataAtualizacao": -1, "last_error_time": 1})
				//findOptions.SetSort(bson.D{{"dataAtualizacao", -1}})
				queryPrevisoes := bson.M{"codINMET": codigoINMET, "dataPrevisao": x.data}
				err = collectionPrevisoes.FindOne(context.TODO(), queryPrevisoes, &queryOptions).Decode(&previsao)
			}
			if err != nil {

				if normais == nil {
					// Erro ao buscar normais
				} else {
					// TODO: calcular o retorno das normais. Talvezes colocar num mapa ou fazer uma função
					//fmt.Println(normais)
					//data :=
					normal := getNormal(x.data.Month().String(), normais)
					//fmt.Println(normal)
					tmin = normal
					tmax = normal
					fonte = "normal"
				}
				// Buscar normais

			} else {
				// *** Pegando dado de PREVISAO ***
				tmin = previsao["temperaturaMinima"].(float64)
				tmax = previsao["temperaturaMaxima"].(float64)
				fonte = "previsao"
				// Debug
				//fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
				//	" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)
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
			//fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
			//	" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)

		}
		// Atualiza o vetor de graus dia
		grausDia := calcularGrausDia(tmin, tmax, TB)
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
	//fmt.Println(gds)

	var st = 0.0
	var qDia = 0.0
	var qPre = 0.0
	var qNor = 0.0
	var qTot = 0.0
	fmt.Println("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
	for i, g := range gds {
		fmt.Printf("Item: %d \t Data: %s \t Graus-dia: %f \t Fonte: %s \t Done: %t\n", i, g.data.Format(layoutISO), g.gd, g.fonte, g.done)
		if g.done {
			st += g.gd
			switch g.fonte {
			case "diario":
				qDia++
			case "previsao":
				qPre++
			case "normal":
				qNor++
			}
			qTot++
		} else {
			fmt.Println("Algum dado nao foi encontrado")
		}

	}
	fmt.Println("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
	fmt.Printf("Intervalo: %s - %s \t Soma Termica: %f (graus dia) \t Diários: %f%% \t Previsões: %f%% \t Normais: %f%% \n", dataInicial, dataFinal, st, ((qDia / qTot) * 100), ((qPre / qTot) * 100), ((qNor / qTot) * 100))

	/*
		TODO: validar dados, calcular graus dia em gds, retornar valor e percentuais
	*/
	pDiario := ((qDia / qTot) * 100)
	pPrevisaoes := ((qPre / qTot) * 100)
	pNormais := ((qNor / qTot) * 100)

	resposta := bson.M{"Soma Termica": st, "Diarios": pDiario, "Previsoes": pPrevisaoes, "Normais": pNormais}

	if gds == nil {
		fmt.Fprintf(w, "Codigo de estação inválido")
	} else {
		json.NewEncoder(w).Encode(resposta)
	}

}

// getDiarios retorna os dados de medições diárias coletados por uma estação meteorológica
func getGrausDiaSudaoProxCorte(w http.ResponseWriter, r *http.Request) {

	// Constantes
	ST_PRIRO_CORTE := 358.00
	ST_OUTROS_CORTES := 281.00

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]
	dataInicial := vars["dataInicial"]
	temperaturaBasal := 10.0
	numeroCortes := vars["numeroCortes"]
	nCortes, err := strconv.ParseInt(numeroCortes, 10, 32)

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
	dataHoje := time.Now()
	dataDiarios := dataHoje.AddDate(0, 0, -2)
	dataPrevisoes := dataHoje.AddDate(0, 0, 14)

	// Cria o slice de datas
	var gds []gd

	// Retornando normais
	var normais bson.M
	queryEstacao := bson.M{"codigoINMET": codigoINMET}
	var estacao bson.M

	err = collectionEstacoes.FindOne(context.TODO(), queryEstacao).Decode(&estacao)
	if err != nil {
		// Erro ao buscar estacoes
		fmt.Println(err)
		fmt.Println("**************************")
		fmt.Println("NORMAIS NAO ENCONTRADAS!!!")
		fmt.Println("**************************")
		return
	} else {
		// Recupera normais
		nomeEstacao := estacao["nomeEstacao"].(string)
		queryNormais := bson.M{"nomeEstacao": nomeEstacao}

		err := collectionNormais.FindOne(context.TODO(), queryNormais).Decode(&normais)
		if err != nil {
			normais = nil
		}
	}

	var st = 0.0
	var dataProximoCorte string
	stAcumulado := 0.0
	dataAtual := dataInicialISO
	for {

		// Se chegou na data atual, retonar a soma termica até aqui
		if dataAtual.Truncate(24 * time.Hour).Equal(dataHoje.Truncate(24 * time.Hour)) {
			//fmt.Println(dataAtual, dataHoje)
			stAcumulado = st
		}

		if (nCortes == 0 && st >= ST_PRIRO_CORTE) || (nCortes > 0 && st >= ST_OUTROS_CORTES) {
			//stAcumulado = st
			dataProximoCorte = dataAtual.Format(layoutISO)
			//fmt.Println("BRACK")
			//fmt.Println(dataProximoCorte)
			//fmt.Println(st)
			//fmt.Println(stAcumulado)
			break
		} else {
			//fmt.Println(st)

			var gdAtual gd

			var tmin, tmax float64
			var fonte string

			var diario bson.M
			var err error = nil
			//Otimizacao
			if dataAtual.After(dataDiarios) {
				err = errors.New("otimizacao")
			} else {
				queryDiarios := bson.M{"codigoINMET": codigoINMET, "dataMedicao": dataAtual}
				err = collectionDiarios.FindOne(context.TODO(), queryDiarios).Decode(&diario)
			}

			if err != nil {
				err = nil
				var previsao bson.M
				//Otimizacao
				if dataAtual.After(dataPrevisoes) {
					err = errors.New("otimizacao")
				} else {
					//fmt.Println("deveria pular")
					//fmt.Println(x.data.Format(layoutISO))
					//fmt.Println(dataHoje.Format(layoutISO))
					queryOptions := options.FindOneOptions{}
					queryOptions.SetSort(bson.M{"dataAtualizacao": -1, "last_error_time": 1})
					//findOptions.SetSort(bson.D{{"dataAtualizacao", -1}})
					queryPrevisoes := bson.M{"codINMET": codigoINMET, "dataPrevisao": dataAtual}
					err = collectionPrevisoes.FindOne(context.TODO(), queryPrevisoes, &queryOptions).Decode(&previsao)
				}
				if err != nil {

					if normais == nil {
						// Erro ao buscar previsoes

						fmt.Println(err)
						fmt.Println("**************************")
						fmt.Println("NORMAIS NAO ENCONTRADAS!!!")
						fmt.Println("**************************")
						fmt.Fprintf(w, "404")

						return

					} else {
						// TODO: calcular o retorno das normais. Talvezes colocar num mapa ou fazer uma função
						//fmt.Println(normais)
						//data :=
						normal := getNormal(dataAtual.Month().String(), normais)
						//fmt.Println(normal)
						tmin = normal
						tmax = normal
						fonte = "normal"
					}
					// Buscar normais

				} else {
					// *** Pegando dado de PREVISAO ***
					tmin = previsao["temperaturaMinima"].(float64)
					tmax = previsao["temperaturaMaxima"].(float64)
					fonte = "previsao"
					// Debug
					//fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
					//	" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)
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
				//fmt.Println("Data: " + x.data.Format(layoutISO) + " | Temperatura Minima: " + fmt.Sprintf("%f", tmin) +
				//	" | Temperatura Maxima: " + fmt.Sprintf("%f", tmax) + " | Fonte: " + fonte)

			}
			// Atualiza o vetor de graus dia
			grausDia := calcularGrausDia(tmin, tmax, temperaturaBasal)
			st += grausDia
			fmt.Println(st)
			if grausDia > 0 {
				gdAtual.gd = grausDia
				gdAtual.fonte = fonte
				gdAtual.done = true
				gdAtual.data = dataAtual
			} else {
				gdAtual.done = false
			}
			gds = append(gds, gdAtual) // Atualiza o slice
			dataAtual = dataAtual.AddDate(0, 0, 1)
		}
	}

	//fmt.Println(gds)
	//fmt.Println("+++++++++++++++++++++++++++++++++++++")
	//fmt.Println(gds)

	var qDia = 0.0
	var qPre = 0.0
	var qNor = 0.0
	var qTot = 0.0

	//fmt.Println("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
	for _, g := range gds {
		//fmt.Printf("Item: %d \t Data: %s \t Graus-dia: %f \t Fonte: %s \t Done: %t\n", i, g.data.Format(layoutISO), g.gd, g.fonte, g.done)

		if g.done {
			//st += g.gd
			switch g.fonte {
			case "diario":
				qDia++
			case "previsao":
				qPre++
			case "normal":
				qNor++
			}
			qTot++
		} else {
			fmt.Println("Algum dado nao foi encontrado")
		}

	}
	//fmt.Println("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
	//fmt.Printf("Intervalo: %s - %s \t Soma Termica: %f (graus dia) \t Diários: %f%% \t Previsões: %f%% \t Normais: %f%% \n", dataInicial, dataFinal, st, ((qDia / qTot) * 100), ((qPre / qTot) * 100), ((qNor / qTot) * 100))

	/*
		TODO: validar dados, calcular graus dia em gds, retornar valor e percentuais
	*/
	pDiario := ((qDia / qTot) * 100)
	pPrevisaoes := ((qPre / qTot) * 100)
	pNormais := ((qNor / qTot) * 100)

	//fmt.Println("FIM")
	//fmt.Println(st)
	//fmt.Println(stAcumulado)
	if stAcumulado == 0 && st > 0 {
		stAcumulado = st
	}
	resposta := bson.M{"proxcorte": dataProximoCorte, "st": stAcumulado, "diario": pDiario, "previsao": pPrevisaoes, "normal": pNormais}

	if gds == nil {
		fmt.Fprintf(w, "Codigo de estação inválido")
	} else {
		json.NewEncoder(w).Encode(resposta)
	}

}

// getPasterjos retorna o número máximo de pastejos de uma área
func getPastejos(w http.ResponseWriter, r *http.Request) {

	// Constantes
	ST_PRIRO_CORTE := 358.00
	ST_OUTROS_CORTES := 281.00

	// Recebe os parametros enviados através da requisição HTTP
	vars := mux.Vars(r)
	codigoINMET := vars["codigoINMET"]
	dataInicial := vars["dataInicial"]
	dataFinal := vars["dataFinal"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collectionNormais := mongoClient.Database("gdsudao").Collection("normais")
	collectionEstacoes := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Conversao das datas para o formato ISO
	layoutISO := "2006-01-02"
	dataInicialISO, _ := time.Parse(layoutISO, dataInicial)
	dataFinalISO, _ := time.Parse(layoutISO, dataFinal)
	dataFinalISO = dataFinalISO.AddDate(0, 0, 1)

	if dataInicialISO.Before(dataFinalISO) {
		// Cria o slice de datas
		var gds []gd
		for currentDate := dataInicialISO; currentDate != dataFinalISO; currentDate = currentDate.AddDate(0, 0, 1) {
			gd := gd{currentDate, 0, "nil", false}
			gds = append(gds, gd)
		}

		// Coletando os dados de normais
		var normais bson.M
		queryEstacao := bson.M{"codigoINMET": codigoINMET}

		var estacao bson.M
		err := collectionEstacoes.FindOne(context.TODO(), queryEstacao).Decode(&estacao)
		if err != nil {
			fmt.Println(err)
		} else {
			nomeEstacao := estacao["nomeEstacao"].(string)
			queryNormais := bson.M{"nomeEstacao": nomeEstacao}
			err := collectionNormais.FindOne(context.TODO(), queryNormais).Decode(&normais)
			if err != nil {
				normais = nil
			}
		}

		if normais == nil {
			resposta := bson.M{"erro": "Codigo de estação inválido"}
			json.NewEncoder(w).Encode(resposta)
		} else {
			var gdTotal = 0.0
			for _, x := range gds {
				normal := getNormal(x.data.Month().String(), normais)
				grausDia := calcularGrausDia(normal, normal, 10.0)
				x.gd = grausDia
				gdTotal += grausDia
				//fmt.Println(x.data.Format(layoutISO) + ": " + fmt.Sprintf("%s", x.gd))
			}

			//fmt.Println(gdTotal)

			var cc = 0.0
			if (gdTotal - ST_PRIRO_CORTE) > 0 {
				cc++
				gdTotal -= ST_PRIRO_CORTE
			} else {
				cc = 0.0
			}

			cc += gdTotal / ST_OUTROS_CORTES
			ccf := fmt.Sprintf("%.2f", cc)

			resposta := bson.M{"Pastejos": ccf}
			if gds == nil {
				fmt.Fprintf(w, "Codigo de estação inválido")
			} else {
				json.NewEncoder(w).Encode(resposta)
			}
		}

	} else {
		resposta := bson.M{"erro": "A data final dever maior que a data inicial"}
		json.NewEncoder(w).Encode(resposta)
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
	myRouter.HandleFunc("/gdsudao/{codigoINMET}/{dataInicial}/{dataFinal}", getGrausDiaSudao).Methods("GET")
	myRouter.HandleFunc("/proximoPastejo/{codigoINMET}/{dataInicial}/{numeroCortes}", getGrausDiaSudaoProxCorte).Methods("GET")
	myRouter.HandleFunc("/grausdia/{codigoINMET}/{temperaturaBasal}/{dataInicial}/{dataFinal}", getGrausDia).Methods("GET")
	myRouter.HandleFunc("/pastejos/{codigoINMET}/{dataInicial}/{dataFinal}", getPastejos).Methods("GET")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

//	main é função Principal do programa
func main() {
	fmt.Println("GDSudão API: on")
	defer fmt.Println("GDSudão API: off")
	handleRequests()
}
