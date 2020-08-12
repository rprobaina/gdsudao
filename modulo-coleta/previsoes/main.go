package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"mongoapi"
	"os"
	"time"
	"utils"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/html/charset"
)

// Forecast tipo de dado que contém os dados retornados pela API do CPTEC
type Forecast struct {
	XMLName     xml.Name   `xml:"cidade"`
	Nome        string     `xml:"nome"`
	UF          string     `xml:"uf"`
	Atualizacao string     `xml:"atualizacao"`
	Previsao    []Previsao `xml:"previsao"`
}

// Previsao tipo de dado que contém os dados de um dia de previsão retornado pela API do CPTEC
type Previsao struct {
	Date   string  `xml:"dia"`
	Tempo  string  `xml:"tempo"`
	Maxima float32 `xml:"maxima"`
	Minima float32 `xml:"minima"`
	Iuv    float32 `xml:"iuv"`
}

// Stations tipo de dado que contem os codigos do CPTEC e do INMET de um estação
type Station struct {
	CodINMET  string `bson:"codigoINMET"`
	CodCPETEC string `bson:"codigoCPTEC"`
}

// getForecastEstendida recupera a previsão do tempo para os póximos 7 dias, de uma estação, da API do CPTEC
func getForecastSevenDays(stationCode string) (Forecast, error) {

	baseURL := "http://servicos.cptec.inpe.br/XML/cidade/7dias/" + stationCode + "/previsao.xml"
	fmt.Println(baseURL)
	utils.Wget(baseURL, "forecast.xml")

	xmlFile, err := os.Open("forecast.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var forecast Forecast
	err = decoder.Decode(&forecast)

	return forecast, err
}

// getForecastEstendida recupera a previsão do tempo estendida, de uma estação, da API do CPTEC
func getForecastEstendida(stationCode string) (Forecast, error) {
	baseURL := "http://servicos.cptec.inpe.br/XML/cidade/" + stationCode + "/estendida.xml"
	fmt.Println(baseURL)
	utils.Wget(baseURL, "forecast.xml")

	xmlFile, err := os.Open("forecast.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var forecast Forecast
	err = decoder.Decode(&forecast)

	return forecast, err
}

// getStations retorna todas as estações inseridas no banco dedados
func getStations() []Station {
	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Consulta estações no banco de dados
	var sources []Station
	filterCursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Atribui as estações retoroanadas a uma lista de estações "sources"
	if err = filterCursor.All(context.TODO(), &sources); err != nil {
		log.Fatal(err)
	}

	// Retorna a lista de estações
	return sources
}

func main() {

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("previsoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Retorna todas as estações presentes no banco de dados
	stations := getStations()

	// Coleta dados de Previsão do tempo para todas as estações
	for _, station := range stations {
		if station.CodCPETEC != "NULL" {

			// Coleta os dados da pevisão para sete dias
			forecast, _ := getForecastSevenDays(station.CodCPETEC)
			for _, previsao := range forecast.Previsao {
				layoutISO := "2006-01-02"
				updateDate, _ := time.Parse(layoutISO, forecast.Atualizacao)
				forecastDate, _ := time.Parse(layoutISO, previsao.Date)
				doc := bson.D{
					{"codINMET", station.CodINMET},
					{"codCPTEC", station.CodCPETEC},
					{"dataPrevisao", forecastDate},
					{"dataAtualizacao", updateDate},
					{"temperaturaMinima", previsao.Minima},
					{"temperaturaMaxima", previsao.Maxima},
					{"iuv", previsao.Iuv},
					{"clima", previsao.Tempo}}
				mongoapi.InsertDocument(*mongoClient, *collection, doc)
			}

			// Coleta os dados da pevisão estendida
			forecast, _ = getForecastEstendida(station.CodCPETEC)
			for _, previsao := range forecast.Previsao {
				layoutISO := "2006-01-02"
				updateDate, _ := time.Parse(layoutISO, forecast.Atualizacao)
				forecastDate, _ := time.Parse(layoutISO, previsao.Date)
				doc := bson.D{
					{"codigoINMET", station.CodINMET},
					{"codigoCPTEC", station.CodCPETEC},
					{"dataPrevisao", forecastDate},
					{"dataAtualizacao", updateDate},
					{"temperaturaMinima", previsao.Minima},
					{"temperaturaMaxima", previsao.Maxima},
					{"iuv", previsao.Iuv},
					{"clima", previsao.Tempo}}
				mongoapi.InsertDocument(*mongoClient, *collection, doc)
			}
		}
	}
}
