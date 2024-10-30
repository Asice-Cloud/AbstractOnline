package session

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func SessionGet(name string, c *gin.Context, key string) any {
	session := sessions.DefaultMany(c, name)
	return session.Get(key)
}

func SessionSet(name string, c *gin.Context, key string, body any) {
	//log.Printf("SessionSet: setting session for %s\n", name) // Add log here
	session := sessions.DefaultMany(c, name)
	if body == nil {
		return
	}
	gob.Register(body)
	session.Set(key, body)
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

func SessionUpdate(name string, c *gin.Context, key string, body any) {
	SessionSet(name, c, key, body)
}

func SessionClear(name string, c *gin.Context) {
	session := sessions.DefaultMany(c, name)
	session.Clear()
}

func SessionDelete(name string, c *gin.Context, key string) {
	session := sessions.DefaultMany(c, name)
	session.Delete(key)
}
