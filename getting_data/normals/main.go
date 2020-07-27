package main

import (
	"fmt"
	"strconv"
	"utils"

	"mongoapi"

	"github.com/extrame/xls"
	"go.mongodb.org/mongo-driver/bson"
)

type normal struct {
	cod  string
	jan  float64
	feb  float64
	mar  float64
	apr  float64
	may  float64
	jun  float64
	jul  float64
	aug  float64
	sep  float64
	oct  float64
	nov  float64
	dec  float64
	year float64
}

func main() {
	// Varibles declaration
	normalsURL := "http://www.inmet.gov.br/webcdp/climatologia/normais2/imagens/normais/planilhas/1981-2010/01%20Temperatura%20M%C3%A9dia%20Compensada%20-%20Bulbo%20Seco%20NCB_1981-2010.xls"

	//dataBaseURI := "mongodb+srv://admin:ricppgcap@cluster0-rmr4a.gcp.mongodb.net/test?retryWrites=true&w=majority"

	dataBaseURI := "mongodb://127.0.0.1:27017"

	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao-test").Collection("normals")
	defer mongoapi.CloseConnection(*mongoClient)

	// Downloading files from INMET
	fmt.Println("Donloading the files...")
	utils.Wget(normalsURL, "normals.xls")
	defer utils.Rm("normals.xls")

	// Parsing stations file
	if stationsFile, err := xls.Open("normals.xls", "utf-8"); err == nil {
		//fmt.Println(stationsFile.GetSheet(0).Row(3).Col(1))
		for i := 4; i <= (int(stationsFile.GetSheet(0).MaxRow)); i++ {
			cod := stationsFile.GetSheet(0).Row(i).Col(0)
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
			normal := normal{cod, jan, feb, mar, apr, may, jun, jul, aug, sep, oct, nov, dec, year}
			fmt.Println(normal)
			doc := bson.D{
				{"cod", normal.cod},
				{"jan", normal.jan},
				{"fev", normal.feb},
				{"mar", normal.mar},
				{"apr", normal.apr},
				{"may", normal.may},
				{"jun", normal.jun},
				{"jul", normal.jul},
				{"aug", normal.aug},
				{"sep", normal.sep},
				{"oct", normal.oct},
				{"nov", normal.nov},
				{"dec", normal.dec},
				{"year", normal.year}}
			if len(cod) > 0 {
				mongoapi.InsertDocument(*mongoClient, *collection, doc)
			}
		}
	}
}
