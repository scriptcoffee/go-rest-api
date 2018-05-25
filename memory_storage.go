package main

type MemoryStore struct {
	id int
	people map[int]Person
}

func (store *MemoryStore) getPeople() ([]Person, error) {
	personList := []Person{}
	for _, value := range store.people {
		personList = append(personList, value)
	}
	return personList, nil
}

func (store *MemoryStore) getPerson(id int) (Person, error) {
	return store.people[id], nil
}

func (store *MemoryStore) createPerson(p Person) error {
	p.Id = store.id
	store.people[p.Id] = p
	store.id++
	return nil
}

func (store *MemoryStore) updatePerson(p Person) error {
	store.people[p.Id] = p
	return nil
}

func (store *MemoryStore) deletePerson(id int) error {
	delete(store.people, id)
	return nil
}
