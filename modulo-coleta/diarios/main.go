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

// Medicao tipo de dados que contem os dados retornados pela API do INMET
type Medicao struct {
	TempMax     string `json:"TEMP_MAX"`
	TempMin     string `json:"TEMP_MIN"`
	TempMed     string `json:"TEMP_MED"`
	DataMedicao string `json:"DT_MEDICAO"`
}

// Station tipo de dado que contem os codigos do CPTEC e do INMET de um estação
type Station struct {
	CodCPTEC string `bson:"codCPTEC"`
	CodINMET string `bson:"codINMET"`
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
	collection := mongoClient.Database("gdsudao").Collection("diarios")
	defer mongoapi.CloseConnection(*mongoClient)

	// Atribui a data do dia anterior à variavel data
	// Data do dia anterior porque o pultimo dado do dia é atualizado no fim de cada dia
	currentTime := time.Now().AddDate(0, 0, -1)
	date := currentTime.Format("2006-01-02")

	// Retorna toda as estações presentes no banco de dados
	stations := getStations()

	// Coleta os dados diário medidos no dia anterior em cada uma das estações cadastradas
	for _, station := range stations {

		// Baixa os dados da API do INMET pada estação da iteração
		baseURL := "https://apitempo.inmet.gov.br/estacao/diaria/" + date + "/" + date + "/" + station.CodINMET
		fmt.Println("Donloading the files...")
		utils.Wget(baseURL, "medicao.json")
		defer utils.Rm("medicao.json")

		// Realiza o parser do dado para uma lista de medições
		medicaoJSONFile, _ := os.Open("medicao.json")
		byteValue, _ := ioutil.ReadAll(medicaoJSONFile)
		var medicoes []Medicao
		json.Unmarshal(byteValue, &medicoes)

		// Valida se o valor existe e adiciona ao banco de dados
		if medicoes[0].TempMin != "" || medicoes[0].TempMax != "" {
			// Realiza as conversões necessárias
			tMax, _ := strconv.ParseFloat(medicoes[0].TempMax, 64)
			tMin, _ := strconv.ParseFloat(medicoes[0].TempMin, 64)
			tMed, _ := strconv.ParseFloat(medicoes[0].TempMed, 64)
			layoutISO := "2006-01-02"
			medicaoDate, _ := time.Parse(layoutISO, medicoes[0].DataMedicao)
			doc := bson.D{
				{"codigoINMET", station.CodINMET},
				{"codigoCPTEC", station.CodCPTEC},
				{"dataMedicao", medicaoDate},
				{"temperaturaMinima", tMin},
				{"temperaturaMaxima", tMax},
				{"temperaturaMedia", tMed},
			}
			mongoapi.InsertDocument(*mongoClient, *collection, doc)
		}
	}
}
