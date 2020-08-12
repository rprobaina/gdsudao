package main

import (
	"fmt"
	"mongoapi"
	"strconv"
	"utils"

	"github.com/extrame/xls"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// URL de dados normais do INMET
	normalsURL := "https://portal.inmet.gov.br/uploads/normais/01-Temperatura-M%C3%A9dia-Compensada-Bulbo-Seco-NCB_1981-2010.xls"

	// Informações de conexão com banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("normais")
	defer mongoapi.CloseConnection(*mongoClient)

	// Baixando os dados de normais climatológicas do INMET
	utils.Wget(normalsURL, "normals.xls")
	defer utils.Rm("normals.xls")

	// Realiza o parser dos dados para a estrutura de dados e insere as informações no banco de dados
	if stationsFile, err := xls.Open("normals.xls", "utf-8"); err == nil {
		for i := 4; i <= (int(stationsFile.GetSheet(0).MaxRow)); i++ {
			cod := stationsFile.GetSheet(0).Row(i).Col(0)
			name := stationsFile.GetSheet(0).Row(i).Col(1)
			jan, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(3), 64)
			feb, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(4), 64)
			mar, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(5), 64)
			apr, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(6), 64)
			may, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(7), 64)
			jun, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(8), 64)
			jul, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(9), 64)
			aug, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(10), 64)
			sep, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(11), 64)
			oct, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(12), 64)
			nov, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(13), 64)
			dec, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(14), 64)
			year, _ := strconv.ParseFloat(stationsFile.GetSheet(0).Row(i).Col(15), 64)
			doc := bson.D{
				{"codigoINMET", cod},
				{"nomeEstacao", name},
				{"normalJaneiro", jan},
				{"normalFevereiro", feb},
				{"normalMarco", mar},
				{"normalAbril", apr},
				{"normalMaio", may},
				{"normalJunho", jun},
				{"normalJulho", jul},
				{"normalAgosto", aug},
				{"normalSetembro", sep},
				{"normalOutubro", oct},
				{"normalNovembro", nov},
				{"normalDezembro", dec},
				{"normalAno", year}}
			mongoapi.InsertDocument(*mongoClient, *collection, doc)
		}
	}
	fmt.Println("Dados de normais climatológicas inseridos com sucesso!")
}
