package main

import (
	"testing"
	"database/sql"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	db *sql.DB
	store Store
}

func (t *TestSuite) SetupSuite() {
	//Connect and set database store
	connString := "dbname=postgres user=web password=<wUA)dXRf6R\\8Z+P"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		t.T().Fatal(err)
	}

	db.Exec("CREATE TABLE IF NOT EXISTS Phonebook (Id serial PRIMARY KEY,Name varchar(255),PhoneNr varchar(255));")

	t.db = db
	t.store = &dbStore{db}
}

func (t *TestSuite) SetupTest() {
	_, err := t.db.Exec("DELETE FROM Phonebook")
	if err != nil {
		t.T().Fatal(err)
	}
	rows, _ := t.db.Query("SELECT SETVAL((SELECT pg_get_serial_sequence('Phonebook', 'id')), 1, false)")
	defer rows.Close()
}

func (t *TestSuite) TearDownSuite() {
	t.db.Close()
}

func TestStoreSuite(t *testing.T) {
	s := new(TestSuite)
	suite.Run(t, s)
}

func (t *TestSuite) insertPeople() {
	t.db.Exec("INSERT INTO Phonebook(name, phonenr) VALUES ($1,$2)", "TestPeter", "8237487324")
	t.db.Exec("INSERT INTO Phonebook(name, phonenr) VALUES ($1,$2)", "TestPaul", "432796597234")
	t.db.Exec("INSERT INTO Phonebook(name, phonenr) VALUES ($1,$2)", "TestAnn", "78264395")
}

func (t *TestSuite) handleMethodCallError(err error, method string) {
	if err != nil {
		t.T().Fatalf(method + " returned error: %s", err)
	}
}

func (t *TestSuite) TestGetPeople() {
	t.insertPeople()

	people, err := t.store.getPeople()

	t.handleMethodCallError(err,"getPeople()")

	peopleCount := len(people)
	if peopleCount != 3 {
		t.T().Fatalf("Expected 3 result rows, got %d", peopleCount)
	}

	firstPerson := people[0]
	if firstPerson.Name != "TestPeter" || firstPerson.PhoneNr != "8237487324" {
		t.T().Fatalf("Expected first person with name 'TestPeter' and phonenr '8237487324', got person with name %s and phonenr %s", firstPerson.Name, firstPerson.PhoneNr)
	}
}

func (t *TestSuite) TestGetPerson() {
	t.insertPeople()

	person, err := t.store.getPerson(2)

	t.handleMethodCallError(err,"getPerson()")
	if person.Name != "TestPaul" || person.PhoneNr != "432796597234" {
		t.T().Fatalf("Expected person with name 'TestPaul' and phonenr '432796597234', got person with name %s and phonenr %s", person.Name, person.PhoneNr)
	}
}

func (t *TestSuite) TestCreatePerson() {
	p := Person{Name:"Charlie", PhoneNr:"8734265034"}
	err := t.store.createPerson(p)

	t.handleMethodCallError(err,"createPerson()")

	rows, err := t.db.Query("SELECT Id, Name, PhoneNr FROM Phonebook WHERE Name = 'Charlie'")

	t.handleMethodCallError(err,"Database query")

	defer rows.Close()

	person := Person{}
	for rows.Next() {
		err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr)
		t.handleMethodCallError(err, "Fetch person data from row")
	}

	if person.Name != "Charlie" || person.PhoneNr != "8734265034" {
		t.T().Fatalf("Expected person with name 'Charlie' and phonenr '8734265034', got person with name %s and phonenr %s", person.Name, person.PhoneNr)
	}
}

func (t *TestSuite) TestUpdatePerson() {
	t.insertPeople()

	p := Person{Id:1, Name:"Bob", PhoneNr:"4398754345"}
	err := t.store.updatePerson(p)

	t.handleMethodCallError(err,"updatePerson()")

	rows, err := t.db.Query("SELECT Id, Name, PhoneNr FROM Phonebook WHERE Name = 'Bob'")

	t.handleMethodCallError(err,"Database query")

	defer rows.Close()

	person := Person{}
	for rows.Next() {
		err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr)
		t.handleMethodCallError(err, "Fetch person data from row")
	}

	if person.Id != 1 || person.Name != "Bob" || person.PhoneNr != "4398754345" {
		t.T().Fatalf("Expected person with name 'Bob' and phonenr '4398754345', got person with name %s and phonenr %s", person.Name, person.PhoneNr)
	}
}

func (t *TestSuite) TestDeletePerson() {
	t.insertPeople()

	err := t.store.deletePerson(1)

	t.handleMethodCallError(err,"deletePerson()")

	rows, err := t.db.Query("SELECT Id, Name, PhoneNr FROM Phonebook")

	t.handleMethodCallError(err,"Database query")

	defer rows.Close()

	people := []Person{}
	for rows.Next() {
		person := Person{}
		err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr);
		t.handleMethodCallError(err, "Fetch person data from row")
		people = append(people, person)
	}
	peopleCount := len(people)
	if peopleCount != 2 {
		t.T().Fatalf("Expected 2 result rows, got %d", peopleCount)
	}

	firstPerson := people[0]
	if firstPerson.Name != "TestPaul" || firstPerson.PhoneNr != "432796597234" {
		t.T().Fatalf("Expected first person with name 'TestPaul' and phonenr '432796597234', got person with name %s and phonenr %s", firstPerson.Name, firstPerson.PhoneNr)
	}
}
