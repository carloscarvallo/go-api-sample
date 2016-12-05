package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/couchbase/gocb"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type Person struct {
	ID        string `json:"id,omitempty"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
}

type N1qlPerson struct {
	Person Person `json:"person"`
}

var bucket *gocb.Bucket

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {

}

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	var person []Person
	query := gocb.NewN1qlQuery("SELECT * FROM `resful-sample` AS person")
	rows, _ := bucket.ExecuteN1qlQuery(query, nil)
	var row N1qlPerson
	for rows.Next(&row) {
		person = append(person, row.Person)
	}
	json.NewEncoder(w).Encode(person)
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	var person Person
	var n1qlParams []interface{}
	_ = json.NewDecoder(req.Body).Decode(&person)
	query := gocb.NewN1qlQuery("INSERT INTO `resful-sample` (KEY, VALUE) values ($1, {'firstname': $2, 'lastname': $3, 'email': $4})")
	n1qlParams = append(n1qlParams, uuid.NewV4().String())
	n1qlParams = append(n1qlParams, person.Firstname)
	n1qlParams = append(n1qlParams, person.Lastname)
	n1qlParams = append(n1qlParams, person.Email)
	_, err := bucket.ExecuteN1qlQuery(query, n1qlParams)
	if err != nil {
		w.WriteHeader(401)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(person)
}

func UpdatePersonEndpoint(w http.ResponseWriter, req *http.Request) {

}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {

}

func main() {
	router := mux.NewRouter()
	cluster, _ := gocb.Connect("couchbase://127.0.0.1")
	bucket, _ = cluster.OpenBucket("resful-sample", "")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/person/{id}", UpdatePersonEndpoint).Methods("POST")
	router.HandleFunc("/person/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
