package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/ddProgerGo/task-kaspi/internal/repository"
	"github.com/ddProgerGo/task-kaspi/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type PersonService struct {
	repo     repository.PersonRepositoryInterface
	validate *validator.Validate
	Logger   *logrus.Logger
	Cache    *redis.Client
}

func NewPersonService(repo repository.PersonRepositoryInterface, logger *logrus.Logger, cache *redis.Client) *PersonService {
	return &PersonService{
		repo:     repo,
		validate: validator.New(),
		Logger:   logger,
		Cache:    cache,
	}
}

func (s *PersonService) SavePerson(person models.Person) error {
	if err := s.validate.Struct(person); err != nil {
		return err
	}

	if _, err := utils.ValidateIIN(person.IIN); err != nil {
		s.Logger.WithError(err).Warn("Invalid IIN format")
		return errors.New("Invalid IIN")
	}
	err := s.repo.SavePerson(person)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to save person")
	} else {
		s.Logger.Info("Person saved successfully: ", person.IIN)
	}
	return err
}

func (s *PersonService) GetPersonByIIN(iin string) (*models.Person, error) {
	ctx := context.Background()

	cached, err := s.Cache.Get(ctx, iin).Result()
	if err == nil {
		var person models.Person
		if err := json.Unmarshal([]byte(cached), &person); err == nil {
			log.Println("Данные загружены из кеша")
			return &person, nil
		}
	}

	person, err := s.repo.GetPersonByIIN(iin)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch person by IIN")
	}

	data, _ := json.Marshal(person)
	s.Cache.Set(ctx, iin, data, 10*time.Minute)

	return person, err
}

func (s *PersonService) GetPeopleByName(name string, page int, limit int) ([]models.Person, error) {
	people, err := s.repo.GetPeopleByName(name, page, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch people by name")
	}
	return people, err
}
