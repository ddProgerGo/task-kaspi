package repository

import (
	"database/sql"
	"errors"

	"github.com/ddProgerGo/task-kaspi/internal/models"
)

type PersonRepository struct {
	DB *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{DB: db}
}

func (r *PersonRepository) SavePerson(person models.Person) error {
	query := `INSERT INTO people (name, iin, phone) VALUES ($1, $2, $3)`
	_, err := r.DB.Exec(query, person.Name, person.IIN, person.Phone)
	return err
}

func (r *PersonRepository) GetPersonByIIN(iin string) (*models.Person, error) {
	query := `SELECT name, iin, phone FROM people WHERE iin = $1`
	row := r.DB.QueryRow(query, iin)
	var person models.Person
	if err := row.Scan(&person.Name, &person.IIN, &person.Phone); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Пользователь не найден")
		}
		return nil, err
	}
	return &person, nil
}

func (r *PersonRepository) GetPeopleByName(namePart string) ([]models.Person, error) {
	query := `SELECT name, iin, phone FROM people WHERE name ILIKE $1`
	rows, err := r.DB.Query(query, "%"+namePart+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		if err := rows.Scan(&person.Name, &person.IIN, &person.Phone); err != nil {
			return nil, err
		}
		people = append(people, person)
	}
	return people, nil
}
