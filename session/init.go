package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InitSession(r *gin.Engine) {
	sessionNames := []string{"user", "admin"}
	//store session in:
	//redis(default)
	//cookie
	store := inRedis()
	r.Use(sessions.SessionsMany(sessionNames, store))
}
