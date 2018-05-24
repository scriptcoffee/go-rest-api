package main

type Store struct {
	id int
	people map[int]Person
}

func (store *Store) getPeople() []Person {
	personList := []Person{}
	for _, value := range store.people {
		personList = append(personList, value)
	}
	return personList
}

func (store *Store) getPerson(id int) Person {
	return store.people[id]
}

func (store *Store) createPerson(p Person) {
	p.Id = store.id
	store.people[p.Id] = p
	store.id++
}

func (store *Store) updatePerson(p Person) {
	store.people[p.Id] = p
}

func (store Store) deletePerson(id int) {
	delete(store.people, id)
}