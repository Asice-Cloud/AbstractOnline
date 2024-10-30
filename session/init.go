package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitSession(r *gin.Engine) {
	//store session in:
	//redis(default)
	//cookie
	store := inRedis()
	r.Use(sessions.SessionsMany(GetAllSession(), store))
	RegisterAll(Table{})
}
