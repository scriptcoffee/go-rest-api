package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"strconv"
	_ "github.com/bmizerany/pq"
	"os"
)

type Person struct {
	Id 		int 	`json:"id" gorm:"primary_key" gorm:"AUTO_INCREMENT"`
	Name 	string 	`json:"name"`
	PhoneNr string 	`json:"phoneNr"`
}

var store Store

func main() {
	//Uncomment to use memory store
	//store = &MemoryStore{0, make(map[int]Person)}

	store = SetupDbStorage()

	//Define routees and methods
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", UpdatePerson).Methods("PUT")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getEnv(name string, defaultVal string) string {
	env := os.Getenv(name)
	if len(env) == 0 {
		env = defaultVal
	}
	return env
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	people := store.getPeople()

	personListBytes, err := json.Marshal(people)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(personListBytes)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	person := Person{}

	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person.Name = r.Form.Get("name")
	person.PhoneNr = r.Form.Get("phoneNr")

	store.createPerson(person)

	w.WriteHeader(http.StatusCreated)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pId := params["id"]

	id, err := strconv.Atoi(pId)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person := store.getPerson(id)

	personBytes, err := json.Marshal(person)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(personBytes)

}

func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	person := Person{}

	err := r.ParseForm()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	params := mux.Vars(r)
	pId := params["id"]

	id, err := strconv.Atoi(pId)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person.Id = id
	person.Name = r.Form.Get("name")
	person.PhoneNr = r.Form.Get("phoneNr")

	store.updatePerson(person)

	w.WriteHeader(http.StatusOK)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	pId := params["id"]

	id, err := strconv.Atoi(pId)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	store.deletePerson(id)

	w.WriteHeader(http.StatusOK)
}