package main

import (
	"context"
	"fmt"
	"log"
	"mongoapi"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

/*
func getAllUsers(w http.ResponseWriter, r *http.Request) {

	URI := "mongodb://localhost:27017"
	client := mong.StartConnection(URI)
	collection := client.Database("golang-test").Collection("users")

	defer mong.CloseConnection(*client)

	users, err := mong.QueryUsers(*client, *collection, bson.D{{}})

	if err != nil {
		fmt.Fprintf(w, "400 - Bad Requestion")
	} else {
		for _, u := range users {
			fmt.Println(u)
		}

		json.NewEncoder(w).Encode(users)
	}

}

func getUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email := vars["email"]

	URI := "mongodb://localhost:27017"
	client := mong.StartConnection(URI)
	collection := client.Database("golang-test").Collection("users")

	defer mong.CloseConnection(*client)

	consulta := bson.D{{"email", email}}
	user, err := mong.QueryUser(*client, *collection, consulta)

	if err != nil {
		fmt.Fprintf(w, "404 - Not Found")
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

func postUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("name")
	userMail := r.FormValue("email")
	userAge := r.FormValue("age")
	fmt.Println(userName, userMail, userAge)
	age, err := strconv.Atoi(userAge)
	if err != nil {
		log.Fatal(err)
	}
	user := mong.User{userName, userMail, age}
	fmt.Println(user)

	URI := "mongodb://localhost:27017"
	client := mong.StartConnection(URI)
	collection := client.Database("golang-test").Collection("users")

	defer mong.CloseConnection(*client)

	mong.InsertUser(*client, *collection, user)

}

func putUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email := vars["email"]

	//	Dados para a atualização
	userName := r.FormValue("name")
	userMail := r.FormValue("email")
	userAge := r.FormValue("age")

	age, err := strconv.Atoi(userAge)
	if err != nil {
		log.Fatal(err)
	}
	user := mong.User{userName, userMail, age}

	URI := "mongodb://localhost:27017"
	client := mong.StartConnection(URI)
	collection := client.Database("golang-test").Collection("users")
	defer mong.CloseConnection(*client)

	consulta := bson.D{{"email", email}}
	atualizacao := bson.D{
		{"$set", bson.D{{"name", user.Name}}},
		{"$set", bson.D{{"email", user.Email}}},
		{"$set", bson.D{{"age", user.Age}}}}

	err = mong.UpdateUser(*client, *collection, consulta, atualizacao)

	if err != nil {
		fmt.Fprintf(w, "400 - Bad Request")
	}
}

//	Delete a user by an email passed in lolcalhost:8081/users/email
func deleteUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	email := vars["email"]

	URI := "mongodb://localhost:27017"
	client := mong.StartConnection(URI)
	collection := client.Database("golang-test").Collection("users")

	defer mong.CloseConnection(*client)

	consulta := bson.D{{"email", email}}

	err := mong.DeleteUser(*client, *collection, consulta)

	if err != nil {
		fmt.Fprintf(w, "400 - Bad Request")
	} else {
		fmt.Fprintf(w, "200 - Ok")
	}
}
*/

// homePage returns a simple message of this API
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a simple REST API writen in Golang that uses MongoDB as database.")
}

// getNearStation retorna os dados da estação meteorológica mais proxima
func getNearStation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	latitude := vars["latitude"]
	longitude := vars["longitude"]

	// Conexão com o banco de dados
	dataBaseURI := "mongodb://127.0.0.1:27017"
	mongoClient := mongoapi.StartConnection(dataBaseURI)
	collection := mongoClient.Database("gdsudao").Collection("estacoes")
	defer mongoapi.CloseConnection(*mongoClient)

	lat, _ := strconv.ParseFloat(latitude, 64)
	lon, _ := strconv.ParseFloat(longitude, 64)

	query := bson.D{
		{"localizacao.coordenadas", bson.D{
			{"$near", bson.D{
				{"$geometry", bson.D{
					{"type", "Point"},
					{"coordinates", bson.A{lat, lon}}},
				}},
			}},
		}}

	//consulta := "{"localizacao.coordenadas": {"$near":{"$geometry":{"type": "Point", "coordinates": [-54.013292, -31.347801]}}}}"

	var resultado bson.D
	err := collection.FindOne(context.TODO(), query).Decode(&resultado)

	fmt.Println(err)
	fmt.Println(resultado)

	// Conexão com o bando de dados
	/*
		dataBaseURI := "mongodb://127.0.0.1:27017"
		mongoClient := mongoapi.StartConnection(dataBaseURI)
		collection := mongoClient.Database("gdsudao").Collection("estacoes")
		defer mongoapi.CloseConnection(*mongoClient)


		consulta := bson.D{{"email", email}}
		user, err := mong.QueryUser(*client, *collection, consulta)
	*/
	//if err != nil {
	fmt.Fprintf(w, "404 - Not Found"+"Lat: "+latitude+"  Lon: "+longitude)
	//} else {
	//json.NewEncoder(w).Encode(user)
	//}
}

//	Trata das requisições (mapeia a requisição para a função adequada)
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/estacao/maisproxima/{latitude}/{longitude}", getNearStation).Methods("GET")
	//myRouter.HandleFunc("/users", getAllUsers).Methods("GET")
	//myRouter.HandleFunc("/users/{email}", getUser).Methods("GET")
	//myRouter.HandleFunc("/users", postUser).Methods("POST")
	//myRouter.HandleFunc("/users/{email}", putUser).Methods("PUT")
	//myRouter.HandleFunc("/users/{email}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

//	Função Principal do programa
func main() {
	fmt.Println("API: on")
	defer fmt.Println("API: off")
	handleRequests()
}
