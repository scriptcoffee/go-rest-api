package main

type Store interface {
	getPeople() ([]Person, error)
	getPerson(id int) (Person, error)
	createPerson(p Person) error
	updatePerson(p Person) error
	deletePerson(id int) error
}