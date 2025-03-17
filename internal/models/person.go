package models

type Person struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required,min=2,max=50"`
	IIN   string `json:"iin" validate:"required,len=12,numeric"`
	Phone string `json:"phone" validate:"required,len=11,numeric"`
}
