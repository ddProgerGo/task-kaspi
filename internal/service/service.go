package service

import "github.com/ddProgerGo/task-kaspi/internal/models"

type PersonServiceInterface interface {
	SavePerson(person models.Person) error
	GetPersonByIIN(iin string) (*models.Person, error)
	GetPeopleByName(name string, page int, limit int) ([]models.Person, error)
}
