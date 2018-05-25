package main

import "database/sql"

type dbStore struct {
	db *sql.DB
}

func (store *dbStore) getPeople() ([]Person, error) {
	rows, err := store.db.Query("SELECT Id, Name, PhoneNr FROM Phonebook")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	people := []Person{}
	for rows.Next() {
		person := Person{}
		if err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr); err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, nil
}

func (store *dbStore) getPerson(id int) (Person, error) {
	rows, err := store.db.Query("SELECT Id, Name, PhoneNr FROM Phonebook WHERE Id=$1", id)
	if err != nil {
		return Person{}, err
	}
	defer rows.Close()

	person := Person{}
	for rows.Next() {
		if err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr); err != nil {
			return Person{}, err
		}
	}
	return person, nil
}

func (store *dbStore) createPerson(p Person) error {
	_, err := store.db.Exec("INSERT INTO Phonebook(name, phonenr) VALUES ($1,$2)", p.Name, p.PhoneNr)
	return err
}

func (store *dbStore) updatePerson(p Person) error {
	_, err := store.db.Query("UPDATE Phonebook SET Name=$1,PhoneNr=$2 WHERE Id=$3", p.Name, p.PhoneNr, p.Id)
	return err
}

func (store *dbStore) deletePerson(id int) error {
	_, err := store.db.Query("DELETE FROM Phonebook WHERE Id=$1", id)
	return err
}
