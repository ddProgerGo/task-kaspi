package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type PersonRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
	Cache  *redis.Client
}

func NewPersonRepository(db *sql.DB, logger *logrus.Logger, cache *redis.Client) *PersonRepository {
	return &PersonRepository{DB: db, Logger: logger}
}

func (r *PersonRepository) SavePerson(person models.Person) error {
	query := `INSERT INTO people (name, iin, phone) VALUES ($1, $2, $3) RETURNING id`
	err := r.DB.QueryRow(query, person.Name, person.IIN, person.Phone).Scan(&person.ID)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("person:%s", person.IIN)
	if err := r.Cache.Del(context.Background(), cacheKey).Err(); err != nil {
		log.Printf("Ошибка очистки кэша для IIN %s: %v", person.IIN, err)
	}

	return nil
}

func (r *PersonRepository) GetPersonByIIN(iin string) (*models.Person, error) {
	query := `SELECT id, name, iin, phone FROM people WHERE iin = $1`
	row := r.DB.QueryRow(query, iin)
	var person models.Person
	if err := row.Scan(&person.ID, &person.Name, &person.IIN, &person.Phone); err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("Person not found with IIN: ", iin)
			return nil, errors.New("Person not found")
		}
		r.Logger.WithError(err).Error("Failed to retrieve person")
		return nil, err
	}

	r.Logger.Info("Person retrieved successfully with IIN: ", iin)
	return &person, nil
}

func (r *PersonRepository) GetPeopleByName(namePart string, page int, limit int) ([]models.Person, error) {
	offset := (page - 1) * limit

	query := `SELECT name, iin, phone FROM people WHERE name ILIKE $1 ORDER BY name ASC LIMIT $2 OFFSET $3`
	rows, err := r.DB.Query(query, "%"+namePart+"%", limit, offset)
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
