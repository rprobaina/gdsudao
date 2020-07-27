package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mongoapi"
	"os"
	"strconv"
	"utils"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/html/charset"
)

// Station get from https://apitempo.inmet.gov.br/estacoes/T
type Station struct {
	CodINMET string `json:"CD_ESTACAO"`
	CodCPTEC string
	Name     string `json:"DC_NOME"`
	Uf       string `json:"SG_ESTADO"`
	Lat      string `json:"VL_LATITUDE"`
	Lon      string `json:"VL_LONGITUDE"`
	Alt      string `json:"VL_ALTITUDE"`
}

type cidades struct {
	XMLName xml.Name `xml:"cidades"`
	Cidade  []Cidade `xml:"cidade"`
}

type Cidade struct {
	Nome string `xml:"nome"`
	Uf   string `xml:"uf"`
	ID   string `xml:"id"`
}

func getCodeCPTEC(cidade string) string {
	baseURLCPETEC := "http://servicos.cptec.inpe.br/XML/listaCidades?city=" + cidade

	// Downloading files from CPTEC
	utils.Wget(baseURLCPETEC, "cidades.xml")
	defer utils.Rm("cidades.xml")

	// Open our xmlFile
	xmlFile, err := os.Open("cidades.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)
	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	// we initialize our Users array
	var cidades cidades
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	//xml.Unmarshal(byteValue, &cidades)

	err = decoder.Decode(&cidades)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	var id string
	fmt.Println(id)
	if len(cidades.Cidade) > 0 {
		id = cidades.Cidade[0].ID
	} else {
		id = "NULL"
	}
	return id
}

func main() {
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

}
