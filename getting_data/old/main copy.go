package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"mongoapi"
	"os"
	"strconv"
	"utils"

	"github.com/extrame/xls"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/html/charset"
)

type station struct {
	cod      string
	codCPTEC string
	name     string
	uf       string
	lat      float64
	lon      float64
	alt      float64
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
	stationsURL := "http://www.inmet.gov.br/webcdp/climatologia/normais2/imagens/normais/planilhas/1981-2010/Esta%C3%A7%C3%B5es%20Normal%20Climato%C3%B3gica%201981-2010.xls"

	//dataBaseURI := "mongodb+srv://admin:ricppgcap@cluster0-rmr4a.gcp.mongodb.net/test?retryWrites=true&w=majority"

	dataBaseURI := "mongodb://127.0.0.1:27017"

	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao-test").Collection("stations")
	defer mongoapi.CloseConnection(*mongoClient)

	// Downloading files from INMET
	fmt.Println("Donloading the files...")
	utils.Wget(stationsURL, "stations.xls")
	defer utils.Rm("stations.xls")

	// Parsing stations file
	if stationsFile, err := xls.Open("stations.xls", "utf-8"); err == nil {
		//fmt.Println(stationsFile.GetSheet(0).Row(3).Col(1))
		for i := 4; i <= (int(stationsFile.GetSheet(0).MaxRow) - 3); i++ {
			cod := stationsFile.GetSheet(0).Row(i).Col(1)
			name := stationsFile.GetSheet(0).Row(i).Col(2)
			codCPTEC := getCodeCPTEC(name)
			uf := stationsFile.GetSheet(0).Row(i).Col(3)
			lat, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(4), 64)
			lon, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(5), 64)
			alt, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(6), 64)
			ws := station{cod, codCPTEC, name, uf, lat, lon, alt}
			fmt.Println(ws)
			doc := bson.D{
				{"cod", ws.cod},
				{"codCPTEC", ws.codCPTEC},
				{"name", ws.name},
				{"uf", ws.uf},
				{"alt", ws.alt},
				{"location", bson.D{
					{"type", "Point"},
					{"coordinates", bson.A{ws.lon, ws.lat}}}}}
			if len(cod) > 0 {
				mongoapi.InsertDocument(*mongoClient, *collection, doc)
			}

		}
	}

}
