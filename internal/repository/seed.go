package repository

import (
	"github.com/google/uuid"
	"github.com/lafetz/assessment/internal/core/domain"
)

func (r *Repository) SeedData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	people := []domain.Person{
		{ID: uuid.New(), Name: "Alice Johnson", Age: 28, Hobbies: []string{"Photography", "Traveling", "Cooking"}},
		{ID: uuid.New(), Name: "Bob Smith", Age: 32, Hobbies: []string{"Reading", "Hiking", "Cycling"}},
		{ID: uuid.New(), Name: "Charlie Brown", Age: 25, Hobbies: []string{"Gaming", "Music", "Drawing"}},
		{ID: uuid.New(), Name: "David Wilson", Age: 30, Hobbies: []string{"Swimming", "Chess", "Writing"}},
		{ID: uuid.New(), Name: "Eve Davis", Age: 27, Hobbies: []string{"Yoga", "Gardening", "Knitting"}},
		{ID: uuid.New(), Name: "Frank Miller", Age: 29, Hobbies: []string{"Surfing", "Video Games", "Rock Climbing"}},
		{ID: uuid.New(), Name: "Grace Lee", Age: 31, Hobbies: []string{"Dance", "Cooking", "Traveling"}},
		{ID: uuid.New(), Name: "Henry Taylor", Age: 26, Hobbies: []string{"Photography", "Blogging", "Reading"}},
		{ID: uuid.New(), Name: "Ivy Martinez", Age: 33, Hobbies: []string{"Fitness", "Traveling", "Painting"}},
		{ID: uuid.New(), Name: "Jack Anderson", Age: 34, Hobbies: []string{"Fishing", "Cooking", "Basketball"}},
		{ID: uuid.New(), Name: "Kimberly Thomas", Age: 22, Hobbies: []string{"Baking", "Reading", "Volunteering"}},
		{ID: uuid.New(), Name: "Liam Jackson", Age: 24, Hobbies: []string{"Coding", "Running", "Music"}},
		{ID: uuid.New(), Name: "Mia White", Age: 35, Hobbies: []string{"Dancing", "Theater", "Crafting"}},
		{ID: uuid.New(), Name: "Noah Harris", Age: 28, Hobbies: []string{"Traveling", "Drawing", "Camping"}},
		{ID: uuid.New(), Name: "Olivia Clark", Age: 31, Hobbies: []string{"Reading", "Yoga", "Traveling"}},
		{ID: uuid.New(), Name: "Paul Lewis", Age: 29, Hobbies: []string{"Video Games", "Football", "Photography"}},
	}

	for _, person := range people {
		r.storage[person.ID] = person
	}
}
