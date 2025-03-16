package handler

import (
	"net/http"

	"github.com/ddProgerGo/task-kaspi/internal/models"
	"github.com/ddProgerGo/task-kaspi/internal/repository"
	"github.com/ddProgerGo/task-kaspi/internal/utils"
	"github.com/gin-gonic/gin"
)

func CheckIIN(c *gin.Context) {
	iin := c.Param("iin")
	info, err := utils.ValidateIIN(iin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"correct": false, "errors": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

func SavePerson(c *gin.Context, repo *repository.PersonRepository) {
	var person models.Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": "Неверный формат данных"})
		return
	}
	if _, err := utils.ValidateIIN(person.IIN); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "errors": err.Error()})
		return
	}
	if err := repo.SavePerson(person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": "Ошибка сохранения"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func GetPersonByIIN(c *gin.Context, repo *repository.PersonRepository) {
	iin := c.Param("iin")
	person, err := repo.GetPersonByIIN(iin)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "errors": "Пользователь не найден"})
		return
	}
	c.JSON(http.StatusOK, person)
}

func GetPeopleByName(c *gin.Context, repo *repository.PersonRepository) {
	name := c.Param("name")
	people, err := repo.GetPeopleByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "errors": "Ошибка поиска"})
		return
	}
	c.JSON(http.StatusOK, people)
}
