package inbound

import (
	"github.com/gin-gonic/gin"
)

type RedirectHandler interface {
	Get(context *gin.Context)
	Post(context *gin.Context)
}
