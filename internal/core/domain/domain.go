package domain

import (
	"github.com/google/uuid"
)

type Person struct {
	ID      uuid.UUID
	Name    string
	Age     int32
	Hobbies []string
}

func NewPerson(
	Name string,
	age int32,
	hobbies []string,
) Person {
	return Person{
		ID:      uuid.New(),
		Name:    Name,
		Age:     age,
		Hobbies: hobbies,
	}
}
