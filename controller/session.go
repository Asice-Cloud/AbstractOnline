package controller

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

type UserSession struct {
	ID       int
	Username string
	Level    int
}

func SessionGet(c *gin.Context, name string) any {
	session := sessions.Default(c)
	return session.Get(name)
}

func SessionSet(c *gin.Context, name string, body any) {
	//log.Printf("SessionSet: setting session for %s\n", name) // Add log here
	session := sessions.Default(c)
	if body == nil {
		return
	}
	gob.Register(body)
	session.Set(name, body)
	// Set session to expire after 30 minutes
	session.Options(sessions.Options{
		MaxAge: 1800, // 30 minutes
	})
	err := session.Save()
	if err != nil {
		log.Printf("Error saving session: %v", err)
	}
	//log.Printf("SessionSet: session set for %s\n", name) // Add log here
}

func SessionUpdate(c *gin.Context, name string, body any) {
	SessionSet(c, name, body)
}

func SessionClear(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
}

func SessionDelete(c *gin.Context, name string) {
	session := sessions.Default(c)
	session.Delete(name)
}
