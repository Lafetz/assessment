package person

import "github.com/google/uuid"

type Person struct {
	ID      uuid.UUID
	Name    string
	age     int32
	hobbies []string
}

func NewPerson(
	Name string,
	age int32,
	hobbies []string,
) Person {
	return Person{
		ID:      uuid.New(),
		age:     age,
		hobbies: hobbies,
	}
}
