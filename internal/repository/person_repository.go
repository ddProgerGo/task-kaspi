package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/ddProgerGo/task-kaspi/pkg/errors"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type PersonRepository struct {
	DB     *sql.DB
	Logger *logrus.Logger
	Cache  *redis.Client
}

func NewPersonRepository(db *sql.DB, logger *logrus.Logger, cache *redis.Client) *PersonRepository {
	return &PersonRepository{DB: db, Logger: logger, Cache: cache}
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
			return nil, errors.ErrNotFound
		}
		r.Logger.WithError(err).Error("Failed to retrieve person")
		return nil, errors.ErrInternalServer
	}

	r.Logger.Info("Person retrieved successfully with IIN: ", iin)
	return &person, nil
}

func (r *PersonRepository) GetPeopleByName(namePart string, page int, limit int) ([]models.Person, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM people WHERE name ILIKE $1`
	if err := r.DB.QueryRow(countQuery, "%"+namePart+"%").Scan(&total); err != nil {
		r.Logger.WithError(err).Error("Failed to get total count of people")
		return nil, 0, err
	}

	query := `SELECT id, name, iin, phone FROM people WHERE name ILIKE $1 ORDER BY name ASC LIMIT $2 OFFSET $3`
	rows, err := r.DB.Query(query, "%"+namePart+"%", limit, offset)
	if err != nil {
		r.Logger.WithError(err).Error("Failed to execute query for people search")
		return nil, 0, err
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		if err := rows.Scan(&person.ID, &person.Name, &person.IIN, &person.Phone); err != nil {
			r.Logger.WithError(err).Error("Failed to scan person row")
			return nil, 0, err
		}
		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		r.Logger.WithError(err).Error("Error iterating through person rows")
		return nil, 0, err
	}

	return people, total, nil
}
