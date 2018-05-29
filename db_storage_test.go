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

	db.Exec("CREATE TABLE IF NOT EXISTS People (Id serial PRIMARY KEY,Name varchar(255),PhoneNr varchar(255));")

	t.db = db
	t.store = SetupDbStorage()
}

func (t *TestSuite) SetupTest() {
	_, err := t.db.Exec("DELETE FROM People")
	if err != nil {
		t.T().Fatal(err)
	}
	rows, _ := t.db.Query("SELECT SETVAL((SELECT pg_get_serial_sequence('People', 'id')), 1, false)")
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
	t.db.Exec("INSERT INTO People(name, phone_nr) VALUES ($1,$2)", "TestPeter", "8237487324")
	t.db.Exec("INSERT INTO People(name, phone_nr) VALUES ($1,$2)", "TestPaul", "432796597234")
	t.db.Exec("INSERT INTO People(name, phone_nr) VALUES ($1,$2)", "TestAnn", "78264395")
}

func (t *TestSuite) handleMethodCallError(err error, method string) {
	if err != nil {
		t.T().Fatalf(method + " returned error: %s", err)
	}
}

func (t *TestSuite) TestGetPeople() {
	t.insertPeople()

	people, _ := t.store.getPeople()

	peopleCount := len(people)
	if peopleCount != 3 {
		t.T().Fatalf("Expected 3 result rows, got %d", peopleCount)
	}

	firstPerson := people[0]
	if firstPerson.Name != "TestPeter" || firstPerson.PhoneNr != "8237487324" {
		t.T().Fatalf("Expected first person with name 'TestPeter' and phone_nr '8237487324', got person with name %s and phone_nr %s", firstPerson.Name, firstPerson.PhoneNr)
	}
}

func (t *TestSuite) TestGetPerson() {
	t.insertPeople()

	person, err := t.store.getPerson(2)
	if err != nil && err.Error() == "Person not found" {
		t.T().Fatalf("Expected person, got error")
	}
	if person.Name != "TestPaul" || person.PhoneNr != "432796597234" {
		t.T().Fatalf("Expected person with name 'TestPaul' and phone_nr '432796597234', got person with name %s and phone_nr %s", person.Name, person.PhoneNr)
	}
}

func (t *TestSuite) TestCreatePerson() {
	p := Person{Name:"Charlie", PhoneNr:"8734265034"}
	t.store.createPerson(p)

	rows, err := t.db.Query("SELECT Id, Name, phone_nr FROM People WHERE Name = 'Charlie'")

	t.handleMethodCallError(err,"Database query")

	defer rows.Close()

	person := Person{}
	for rows.Next() {
		err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr)
		t.handleMethodCallError(err, "Fetch person data from row")
	}

	if person.Name != "Charlie" || person.PhoneNr != "8734265034" {
		t.T().Fatalf("Expected person with name 'Charlie' and phone_nr '8734265034', got person with name %s and phone_nr %s", person.Name, person.PhoneNr)
	}
}

func (t *TestSuite) TestUpdatePerson() {
	t.insertPeople()

	p := Person{Id:1, Name:"Bob", PhoneNr:"4398754345"}
	t.store.updatePerson(p)

	rows, err := t.db.Query("SELECT Id, Name, phone_nr FROM People WHERE Name = 'Bob'")

	t.handleMethodCallError(err,"Database query")

	defer rows.Close()

	person := Person{}
	for rows.Next() {
		err := rows.Scan(&person.Id, &person.Name, &person.PhoneNr)
		t.handleMethodCallError(err, "Fetch person data from row")
	}

	if person.Id != 1 || person.Name != "Bob" || person.PhoneNr != "4398754345" {
		t.T().Fatalf("Expected person with name 'Bob' and phone_nr '4398754345', got person with name %s and phone_nr %s", person.Name, person.PhoneNr)
	}
}

func (t *TestSuite) TestDeletePerson() {
	t.insertPeople()

	t.store.deletePerson(1)

	rows, err := t.db.Query("SELECT Id, Name, phone_nr FROM People")

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
		t.T().Fatalf("Expected first person with name 'TestPaul' and phone_nr '432796597234', got person with name %s and phone_nr %s", firstPerson.Name, firstPerson.PhoneNr)
	}
}
