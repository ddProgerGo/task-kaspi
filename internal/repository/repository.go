package repository

import "github.com/ddProgerGo/task-kaspi/internal/models"

type PersonRepositoryInterface interface {
	SavePerson(person models.Person) error
	GetPersonByIIN(iin string) (*models.Person, error)
	GetPeopleByName(namePart string, page int, limit int) ([]models.Person, int, error)
}
