package main

import (
	"github.com/jinzhu/gorm"
	"errors"
)

type dbStore struct {
	db *gorm.DB
}

func (store *dbStore) getPeople() ([]Person, error) {
	person := []Person{}
	err := store.db.Find(&person).Error
	return person, err
}

func (store *dbStore) getPerson(id int) (Person, error) {
	person := Person{}
	notFound := store.db.First(&person, id).RecordNotFound()
	if notFound {
		return person, errors.New("Person not found")
	}
	return person, nil
}

func (store *dbStore) createPerson(p Person) error {
	return store.db.Create(&p).Error
}

func (store *dbStore) updatePerson(p Person) error {
	return store.db.Save(&p).Error
}

func (store *dbStore) deletePerson(id int) error {
	person := Person{Id:id}
	return store.db.Delete(&person).Error
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
