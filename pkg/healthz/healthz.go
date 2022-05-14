package healthz

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "health",
	})
}
