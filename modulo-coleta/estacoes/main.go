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

// Station representa os dados coletados na API do INMET
type Station struct {
	CodINMET string `json:"CD_ESTACAO"`
	CodCPTEC string
	Name     string `json:"DC_NOME"`
	Tipo     string `json:"TP_ESTACAO"`
	Uf       string `json:"SG_ESTADO"`
	Lat      string `json:"VL_LATITUDE"`
	Lon      string `json:"VL_LONGITUDE"`
	Alt      string `json:"VL_ALTITUDE"`
}

// cidades representa uma lista de cidades
type cidades struct {
	XMLName xml.Name `xml:"cidades"`
	Cidade  []Cidade `xml:"cidade"`
}

// Cidade representa os dados coletados na API do CPTEC
type Cidade struct {
	Nome string `xml:"nome"`
	Uf   string `xml:"uf"`
	ID   string `xml:"id"`
}

// getCodeCPTEC retorna o codigo de estação utilizado pelo CPTEC
func getCodeCPTEC(cidade string) string {

	// Baixa os dados da API do CPTEC
	baseURLCPETEC := "http://servicos.cptec.inpe.br/XML/listaCidades?city=" + cidade
	utils.Wget(baseURLCPETEC, "cidades.xml")
	defer utils.Rm("cidades.xml")

	// Decodifica os dados em um ByteArray
	xmlFile, err := os.Open("cidades.xml")
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	// Faz o parser dos dados para uma lista de cidades
	var cidades cidades
	err = decoder.Decode(&cidades)

	// Valida se o ide da estação foi encontrado e retorna o resultado
	var id string
	if len(cidades.Cidade) > 0 {
		id = cidades.Cidade[0].ID
	} else {
		id = "NULL"
	}
	return id
}

func getStations(stationsURL string) {
	// Conexão com o bando de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	// Baixa dos dados do from INMET
	utils.Wget(stationsURL, "estacoes.json")
	defer utils.Rm("estacoes.json")

	// Faz o parser do arquivo para um ByteArray
	stationsJSONFile, _ := os.Open("estacoes.json")
	byteValue, _ := ioutil.ReadAll(stationsJSONFile)

	// Trnsfere os dados to ByteArray para uma lista de estações
	var stations []Station
	json.Unmarshal(byteValue, &stations)

	// Formata e insere os dados no banco de dados
	for _, station := range stations {
		codCPTEC := getCodeCPTEC(station.Name)
		alt, _ := strconv.ParseFloat(station.Alt, 64)
		lat, _ := strconv.ParseFloat(station.Lat, 64)
		lon, _ := strconv.ParseFloat(station.Lon, 64)
		doc := bson.D{
			{"codigoINMET", station.CodINMET},
			{"codigoCPTEC", codCPTEC},
			{"nomeEstacao", station.Name},
			{"tipoEstacao", station.Tipo},
			{"uf", station.Uf},
			{"altitude", alt},
			{"localizacao", bson.D{
				{"type", "Point"},
				{"coordenadas", bson.A{lon, lat}}}}}
		mongoapi.InsertDocument(*mongoClient, *collection, doc)
	}
}

func main() {
	// URL dos dados das estações
	manualStationsURL := "https://apitempo.inmet.gov.br/estacoes/M"
	automaticStationsURL := "https://apitempo.inmet.gov.br/estacoes/T"

	// Insere estações manuais
	getStations(manualStationsURL)

	// Insere estações automaticas
	getStations(automaticStationsURL)
}
