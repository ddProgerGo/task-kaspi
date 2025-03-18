package middleware

import (
	"net/http"

	"github.com/ddProgerGo/task-kaspi/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorHandlingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.WithField("panic", r).Error("Panic recovered in middleware")
				c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServer.Message})
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.WithError(err).Error("Request error")

				if appErr, ok := err.Err.(*errors.AppError); ok {

					if appErr.IsDefault {
						c.JSON(appErr.Code, gin.H{"success": false, "error": appErr.Message})
						c.Abort()
						return
					} else {
						c.JSON(appErr.Code, gin.H{"correct": false, "error": appErr.Message})
						c.Abort()
						return
					}
				}
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServer.Message})
			c.Abort()
		}
	}
}
