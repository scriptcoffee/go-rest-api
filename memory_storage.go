package main

type MemoryStore struct {
	id int
	people map[int]Person
}

func (store *MemoryStore) getPeople() []Person {
	personList := []Person{}
	for _, value := range store.people {
		personList = append(personList, value)
	}
	return personList
}

func (store *MemoryStore) getPerson(id int) Person {
	return store.people[id]
}

func (store *MemoryStore) createPerson(p Person) {
	p.Id = store.id
	store.people[p.Id] = p
	store.id++
}

func (store *MemoryStore) updatePerson(p Person) {
	store.people[p.Id] = p
}

func (store *MemoryStore) deletePerson(id int) {
	delete(store.people, id)
}
