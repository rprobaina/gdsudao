package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/html/charset"
)

// the user struct, this contains our
// Type attribute, our user's name and
// a social struct which will contain all
// our social links
type Cidade struct {
	XMLName     xml.Name   `xml:"cidade"`
	Nome        string     `xml:"nome"`
	UF          string     `xml:"uf"`
	Atualizacao string     `xml:"atualizacao"`
	Previsao    []Previsao `xml:"previsao"`
}

type Previsao struct {
	Dia    string `xml:"dia"`
	Tempo  string `xml:"tempo"`
	Maxima string `xml:"maxima"`
	Minima string `xml:"minima"`
	Iuv    string `xml:"iuv"`
}

func main() {

	// Open our xmlFile
	xmlFile, err := os.Open("teste.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	reader := bytes.NewReader(byteValue)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	// we initialize our Users array
	var cidades Cidade
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	//xml.Unmarshal(byteValue, &cidades)

	err = decoder.Decode(&cidades)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	fmt.Println(cidades)
	fmt.Println(cidades.Atualizacao)
	fmt.Println(cidades.Previsao[1].Dia)
}
