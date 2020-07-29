package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mongoapi"
	"os"
	"strconv"
	"time"
	"utils"

	"go.mongodb.org/mongo-driver/bson"
)

// Station get from https://apitempo.inmet.gov.br/estacoes/T
type Medicao struct {
	TempMax     string `json:"TEMP_MAX"`
	TempMin     string `json:"TEMP_MIN"`
	TempMed     string `json:"TEMP_MED"`
	DataMedicao string `json:"DT_MEDICAO"`
}

type Station struct {
	//CodCPETEC string `bson:"codCPTEC"`
	CodINMET string `bson:"codINMET"`
}

func getStations() []Station {
	dataBaseURI := "mongodb://127.0.0.1:27017"

	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao-test").Collection("stations")
	defer mongoapi.CloseConnection(*mongoClient)

	var stations []Station
	filterCursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = filterCursor.All(context.TODO(), &stations); err != nil {
		log.Fatal(err)
	}

	return stations
}

func main() {

	//dataBaseURI := "mongodb+srv://admin:ricppgcap@cluster0-rmr4a.gcp.mongodb.net/test?retryWrites=true&w=majority"
	dataBaseURI := "mongodb://127.0.0.1:27017"

	// Open dadabase connection
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao-test").Collection("real")
	defer mongoapi.CloseConnection(*mongoClient)

	// Return the yesterday data in the format YYYY-MM-DD
	currentTime := time.Now().AddDate(0, 0, -1)
	date := currentTime.Format("2006-01-02")

	fmt.Println(date)

	// https://apitempo.inmet.gov.br/estacao/diaria/2019-10-01/2019-10-31/A301

	stations := getStations()

	for _, station := range stations {
		baseURL := "https://apitempo.inmet.gov.br/estacao/diaria/" + date + "/" + date + "/" + station.CodINMET

		fmt.Println(baseURL)
		// Downloading files from INMET
		fmt.Println("Donloading the files...")
		utils.Wget(baseURL, "medicao.json")
		defer utils.Rm("medicao.json")

		// Parsing stations file

		medicaoJSONFile, _ := os.Open("medicao.json")
		//defer os.Close(stationsJsonFile)

		// read our opened jsonFile as a byte array.
		byteValue, _ := ioutil.ReadAll(medicaoJSONFile)

		var medicoes []Medicao
		json.Unmarshal(byteValue, &medicoes)

		//fmt.Println(medicoes)

		if medicoes[0].TempMin != "" {
			tMax, _ := strconv.ParseFloat(medicoes[0].TempMax, 64)
			tMin, _ := strconv.ParseFloat(medicoes[0].TempMin, 64)
			tMed, _ := strconv.ParseFloat(medicoes[0].TempMed, 64)
			layoutISO := "2006-01-02"
			medicaoDate, _ := time.Parse(layoutISO, medicoes[0].DataMedicao)
			doc := bson.D{
				{"codINMET", station.CodINMET},
				{"date", medicaoDate},
				{"temperature_min", tMin},
				{"temperature_max", tMax},
				{"temperature_med", tMed},
			}

			mongoapi.InsertDocument(*mongoClient, *collection, doc)
			//fmt.Println(doc)
		}

	}

	/*
		// Varibles declaration
		stationsURL := "https://apitempo.inmet.gov.br/estacoes/T"

		//dataBaseURI := "mongodb+srv://admin:ricppgcap@cluster0-rmr4a.gcp.mongodb.net/test?retryWrites=true&w=majority"
		dataBaseURI := "mongodb://127.0.0.1:27017"

		// Open dadabase connection
		mongoClient := mongoapi.StartConnection(dataBaseURI)
		collection := mongoClient.Database("gdsudao-test").Collection("stations")
		defer mongoapi.CloseConnection(*mongoClient)

		// Downloading files from INMET
		fmt.Println("Donloading the files...")
		utils.Wget(stationsURL, "stations.json")
		defer utils.Rm("stations.json")

		// Parsing stations file

		stationsJSONFile, _ := os.Open("stations.json")
		//defer os.Close(stationsJsonFile)

		// read our opened jsonFile as a byte array.
		byteValue, _ := ioutil.ReadAll(stationsJSONFile)

		var stations []Station
		json.Unmarshal(byteValue, &stations)

		for _, station := range stations {
			codCPTEC := getCodeCPTEC(station.Name)
			alt, _ := strconv.ParseFloat(station.Alt, 64)
			lat, _ := strconv.ParseFloat(station.Lat, 64)
			lon, _ := strconv.ParseFloat(station.Lon, 64)
			doc := bson.D{
				{"codINMET", station.CodINMET},
				{"codCPTEC", codCPTEC},
				{"name", station.Name},
				{"uf", station.Uf},
				{"alt", alt},
				{"location", bson.D{
					{"type", "Point"},
					{"coordinates", bson.A{lon, lat}}}}}

			mongoapi.InsertDocument(*mongoClient, *collection, doc)
		}
	*/

}
