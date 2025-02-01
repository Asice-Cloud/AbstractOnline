package user_module

import (
	"Abstract/config"

	"github.com/gin-gonic/gin"
)

func Index(context *gin.Context) {
	config.Lg.Debug("Index")
}
