package main

import (
	"github.com/jinzhu/gorm"
)

type dbStore struct {
	db *gorm.DB
}

func (store *dbStore) getPeople() []Person {
	person := []Person{}
	store.db.Find(&person)
	return person
}

func (store *dbStore) getPerson(id int) Person {
	person := Person{}
	store.db.First(&person, id)
	return person
}

func (store *dbStore) createPerson(p Person)  {
	store.db.Create(&p)
}

func (store *dbStore) updatePerson(p Person)  {
	store.db.Save(&p)
}

func (store *dbStore) deletePerson(id int) {
	person := Person{Id:id}
	store.db.Delete(&person)
}

func SetupDbStorage() Store {
	//Connect and set database store
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	dbname := getEnv("DB_NAME", "postgres")
	user := getEnv("DB_USER", "web")
	password := getEnv("DB_PASSWORD", "<wUA)dXRf6R\\8Z+P")
	sslMode := getEnv("DB_SSL_MODE", "disable")
	connString := "host=" + host + " port=" + port + " dbname=" + dbname + " user=" + user + " password=" + password + " sslmode=" + sslMode
	db, err := gorm.Open("postgres", connString)

	if err != nil {
		panic(err)
	}

	if !db.HasTable(&Person{}) {
		db.CreateTable(&Person{})
	}

	return &dbStore{db}
}
