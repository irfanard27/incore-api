package httpresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"message": message, "data": data})
}

func ErrorResponse(c *gin.Context, err error, code int) {
	c.JSON(code, gin.H{"error": err.Error()})
}
