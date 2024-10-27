package dto

type CreatePerson struct {
	Name    string   `json:"name" validate:"required"`
	Age     int32    `json:"age" validate:"required,gte=0,lte=120"`
	Hobbies []string `json:"hobbies" validate:"required,dive,required"`
}
type UpdatePerson struct {
	Name    string   `json:"name" validate:"required"`
	Age     int32    `json:"age" validate:"required,gte=0,lte=120"`
	Hobbies []string `json:"hobbies" validate:"required,dive,required"`
}
