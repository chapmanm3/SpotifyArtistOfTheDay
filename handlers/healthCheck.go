package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h DBHandlerService) GetHealthCheck(c *gin.Context) {
	fmt.Println("In Health Check")
  c.IndentedJSON(http.StatusOK, "Looking Good")
}
