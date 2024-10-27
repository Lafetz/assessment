package dto

import (
	"github.com/google/uuid"
	"github.com/lafetz/assessment/internal/core/domain"
)

type JSONPerson struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Age     int32     `json:"age"`
	Hobbies []string  `json:"hobbies"`
}

func ConvertToJSONPerson(p domain.Person) JSONPerson {
	return JSONPerson{
		ID:      p.ID,
		Name:    p.Name,
		Age:     p.Age,
		Hobbies: p.Hobbies,
	}
}
func ConvertToJSONPersonArray(persons []domain.Person) []JSONPerson {
	jsonPersons := make([]JSONPerson, len(persons))
	for i, p := range persons {
		jsonPersons[i] = ConvertToJSONPerson(p)
	}
	return jsonPersons
}

type JSONMetadata struct {
	CurrentPage  int32 `json:"currentPage"`
	PageSize     int32 `json:"pageSize"`
	FirstPage    int32 `json:"firstPage"`
	LastPage     int32 `json:"lastPage"`
	TotalRecords int32 `json:"totalRecords"`
}

func ConvertToJSONMetadata(meta domain.Metadata) JSONMetadata {
	return JSONMetadata{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		FirstPage:    meta.FirstPage,
		LastPage:     meta.LastPage,
		TotalRecords: meta.TotalRecords,
	}
}

type GetPersonsResponse struct {
	Meta    JSONMetadata `json:"meta"`
	Persons []JSONPerson `json:"persons"`
}

func ConvertToGetPersonsResponse(persons []domain.Person, meta domain.Metadata) GetPersonsResponse {
	return GetPersonsResponse{
		Meta:    ConvertToJSONMetadata(meta),
		Persons: ConvertToJSONPersonArray(persons),
	}
}
