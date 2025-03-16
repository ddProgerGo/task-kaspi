package handler

import (
	"database/sql"

	"github.com/ddProgerGo/task-kaspi/internal/repository"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	repo := repository.NewPersonRepository(db)
	r.GET("/iin_check/:iin", CheckIIN)
	r.POST("/people/info", func(c *gin.Context) { SavePerson(c, repo) })
	r.GET("/people/info/iin/:iin", func(c *gin.Context) { GetPersonByIIN(c, repo) })
	r.GET("/people/info/phone/:name", func(c *gin.Context) { GetPeopleByName(c, repo) })
}
