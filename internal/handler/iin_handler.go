package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
