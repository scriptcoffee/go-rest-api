package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
	"fmt"
	"strconv"
	_ "github.com/bmizerany/pq"
	"database/sql"
	"os"
)

type Person struct {
	Id 		int 	`json:"id"`
	Name 	string 	`json:"name"`
	PhoneNr string 	`json:"phoneNr"`
}

var store Store

func main() {
	//Uncomment to use memory store
	//store = &MemoryStore{0, make(map[int]Person)}

	//Connect and set database store
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "postgres")
	user := getEnv("DB_USER", "web")
	password := getEnv("DB_PASSWORD", "<wUA)dXRf6R\\8Z+P")
	sslMode := getEnv("DB_SSL_MODE", "disable")
	connString := "host=" + host + " port=" + port + " dbname=" + dbname + " user=" + user + " password=" + password + " sslmode=" + sslMode
	db, err := sql.Open("postgres", connString)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	store = &dbStore{db}


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
	people, err := store.getPeople()

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	err = store.createPerson(person)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	person, err := store.getPerson(id)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	err = store.updatePerson(person)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

	err = store.deletePerson(id)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}