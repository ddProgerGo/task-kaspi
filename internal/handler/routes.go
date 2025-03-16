package handler

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/iin_check/:iin", CheckIIN)
	r.POST("/people/info", SavePerson)
	r.GET("/people/info/iin/:iin", GetPersonByIIN)
	r.GET("/people/info/phone/:name", GetPeopleByName)
}
