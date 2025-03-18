package handler

import (
	"net/http"
	"strconv"

	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/ddProgerGo/task-kaspi/internal/service"
	"github.com/ddProgerGo/task-kaspi/internal/utils"
	"github.com/ddProgerGo/task-kaspi/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type PersonHandler struct {
	service service.PersonServiceInterface
	Logger  *logrus.Logger
}

func NewPersonHandler(service service.PersonServiceInterface, logger *logrus.Logger) *PersonHandler {
	return &PersonHandler{service: service, Logger: logger}
}

// CheckIIN godoc
// @Summary     Validate IIN
// @Description Checks if the provided IIN is valid
// @Tags        IIN
// @Accept      json
// @Produce     json
// @Param       iin  path  string  true  "IIN number"
// @Success     200  {object}  map[string]interface{}
// @Failure     400  {object}  map[string]string
// @Router      /check-iin/{iin} [get]
func (h *PersonHandler) CheckIIN(c *gin.Context) {
	iin := c.Param("iin")

	info, err := utils.ValidateIIN(iin)
	if err != nil {
		h.Logger.WithError(err).Warn("Invalid IIN check")
		c.Error(err)
		return
	}

	h.Logger.Info("IIN validation successful: ", iin)
	c.JSON(http.StatusOK, info)
}

// SavePerson godoc
// @Summary     Save a person
// @Description Saves a new person to the database
// @Tags        Person
// @Accept      json
// @Produce     json
// @Param       person  body  models.Person  true  "Person data"
// @Success     200  {object}  map[string]bool
// @Failure     400  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /save-person [post]
func (h *PersonHandler) SavePerson(c *gin.Context) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		h.Logger.WithError(err).Warn("Invalid request format")
		c.Error(errors.ErrBadRequest)
		return
	}

	if err := h.service.SavePerson(person); err != nil {
		h.Logger.WithError(err).Error("Failed to save person")

		if appErr, ok := err.(*errors.AppError); ok {
			c.Error(&errors.AppError{Code: appErr.Code, Message: appErr.Message, IsDefault: true})
		} else {
			c.Error(&errors.AppError{Code: http.StatusInternalServerError, Message: err.Error(), IsDefault: true})
		}
		return
	}

	h.Logger.Info("Person saved successfully: ", person.IIN)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPersonByIIN godoc
// @Summary     Get person by IIN
// @Description Retrieves person details by IIN
// @Tags        Person
// @Accept      json
// @Produce     json
// @Param       iin  path  string  true  "IIN number"
// @Success     200  {object}  models.Person
// @Failure     404  {object}  map[string]string
// @Router      /get-person/{iin} [get]
func (h *PersonHandler) GetPersonByIIN(c *gin.Context) {
	iin := c.Param("iin")

	person, err := h.service.GetPersonByIIN(iin)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.Error(&errors.AppError{Code: appErr.Code, Message: appErr.Message, IsDefault: true})
		} else {
			c.Error(errors.ErrNotFound)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": person})
}

// GetPeopleByName godoc
// @Summary     Get people by name with pagination
// @Description Retrieves a paginated list of people matching the provided name
// @Tags        Person
// @Accept      json
// @Produce     json
// @Param       name   path      string  true  "Person name"
// @Param       page   query     int     false "Page number" default(1)
// @Param       limit  query     int     false "Results per page" default(10)
// @Success     200    {array}   models.Person
// @Failure     500    {object}  map[string]string
// @Router      /get-people/{name} [get]
func (h *PersonHandler) GetPeopleByName(c *gin.Context) {
	name := c.Param("name")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		h.Logger.Warn("Invalid page number")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		h.Logger.Warn("Invalid limit number")
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": "Invalid limit number"})
		return
	}

	people, total, err := h.service.GetPeopleByName(name, page, limit)
	if err != nil {
		h.Logger.WithError(err).Error("Error searching people")
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": "Error searching people"})
		return
	}

	if len(people) <= 0 {
		h.Logger.Warn("No people found for name:", name)
		people = []models.Person{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    people,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}
