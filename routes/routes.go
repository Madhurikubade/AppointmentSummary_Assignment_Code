package routes

import (
	"AppointmentSummary_Assignment_Code/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(server *gin.Engine) {
	server.GET(config.AppConfig.Prefix+"/ping", PingResponse)

}

func PingResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "PONG"})
}
