package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/ddProgerGo/task-kaspi/internal/repository"
	"github.com/ddProgerGo/task-kaspi/internal/utils"
	"github.com/ddProgerGo/task-kaspi/pkg/errors"
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
		return errors.ErrBadRequest
	}

	if _, err := utils.ValidateIIN(person.IIN); err != nil {
		s.Logger.WithError(err).Warn("Invalid IIN format")
		return err
	}

	err := s.repo.SavePerson(person)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to save person: ", err)
		return err
	}

	s.Logger.Info("Person saved successfully: ", person.IIN)
	return nil
}

func (s *PersonService) GetPersonByIIN(iin string) (*models.Person, error) {
	ctx := context.Background()

	if _, err := utils.ValidateIIN(iin); err != nil {
		s.Logger.WithError(err).Warn("Invalid IIN format")
		return nil, err
	}

	cached, err := s.Cache.Get(ctx, iin).Result()
	if err == nil {
		var person models.Person
		if err := json.Unmarshal([]byte(cached), &person); err == nil {
			s.Logger.Info("Data loaded from cache for IIN: ", iin)
			return &person, nil
		}
	}

	person, err := s.repo.GetPersonByIIN(iin)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch person by IIN")
		return nil, err
	}

	data, err := json.Marshal(person)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to serialize person data for caching")
		return nil, errors.ErrInternalServer
	}

	if err := s.Cache.Set(ctx, iin, data, 10*time.Minute).Err(); err != nil {
		s.Logger.WithError(err).Error("Failed to cache person data")
	}
	return person, nil
}

func (s *PersonService) GetPeopleByName(name string, page int, limit int) ([]models.Person, int, error) {
	people, total, err := s.repo.GetPeopleByName(name, page, limit)
	if err != nil {
		s.Logger.WithError(err).Error("Failed to fetch people by name")
	}
	return people, total, err
}
