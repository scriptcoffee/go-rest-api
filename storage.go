package main

type Store interface {
	getPeople() []Person
	getPerson(id int) Person
	createPerson(p Person)
	updatePerson(p Person)
	deletePerson(id int)
}