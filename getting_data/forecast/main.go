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

type Forecast struct {
	XMLName     xml.Name   `xml:"cidade"`
	Nome        string     `xml:"nome"`
	UF          string     `xml:"uf"`
	Atualizacao string     `xml:"atualizacao"`
	Previsao    []Previsao `xml:"previsao"`
}

type Previsao struct {
	Date   string  `xml:"dia"`
	Tempo  string  `xml:"tempo"`
	Maxima float32 `xml:"maxima"`
	Minima float32 `xml:"minima"`
	Iuv    float32 `xml:"iuv"`
}

type Station struct {
	StationId string `bson:"_id"`
	CodCPETEC string `bson:"codCPTEC"`
}

func getForecastForFourDays(stationCode string) (Forecast, error) {
	// Base URL from cptec services
	baseURL := "http://servicos.cptec.inpe.br/XML/cidade/" + stationCode + "/previsao.xml"

	// Downloading files from CPTEC
	utils.Wget(baseURL, "forecast.xml")

	xmlFile, err := os.Open("forecast.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)
	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var forecast Forecast
	err = decoder.Decode(&forecast)
	return forecast, err
}

func getForecastForSevenDays(stationCode string) (Forecast, error) {
	// Base URL from cptec services
	baseURL := "http://servicos.cptec.inpe.br/XML/cidade/7dias/" + stationCode + "/previsao.xml"

	// Downloading files from CPTEC
	utils.Wget(baseURL, "forecast.xml")

	xmlFile, err := os.Open("forecast.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)
	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var forecast Forecast
	err = decoder.Decode(&forecast)

	return forecast, err
}

func getStations() []Station {
	dataBaseURI := "mongodb://127.0.0.1:27017"

	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao-test").Collection("stations")
	defer mongoapi.CloseConnection(*mongoClient)

	var sources []Station
	filterCursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = filterCursor.All(context.TODO(), &sources); err != nil {
		log.Fatal(err)
	}

	return sources
}

func main() {
	dataBaseURI := "mongodb://127.0.0.1:27017"

	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao-test").Collection("forecasts")
	defer mongoapi.CloseConnection(*mongoClient)

	sources := getForecastSources()

	//fmt.Println(sources)

	for _, source := range sources {
		if source.CodCPETEC != "NULL" {
			forecast, _ := getForecastForSevenDays(source.CodCPETEC)
			fmt.Println("ID: " + source.StationId)
			fmt.Println("Source: " + source.CodCPETEC)
			//fmt.Println(forecast)
			layoutISO := "2006-01-02"
			updateDate, _ := time.Parse(layoutISO, forecast.Atualizacao)
			doc := bson.D{
				{"stationId", source.StationId},
				{"updateDate", updateDate},
				{"forecasts", forecast.Previsao}}

			mongoapi.InsertDocument(*mongoClient, *collection, doc)
			fmt.Println("inserido")
			fmt.Println("==========================================")
		}
	}
	//forecastFour, _ := getForecastForFourDays("694")
	//fmt.Println(forecastFour)
	forecastSeven, _ := getForecastForSevenDays("694")
	fmt.Println(forecastSeven)
	//fmt.Println(forecast.Atualizacao)
	//fmt.Println(forecast.Previsao[1].Dia)

}
