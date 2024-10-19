package session

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
	"reflect"
	"strings"
)

var session_name_list []string

func inRedis() sessions.Store {
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	return store
}

func inCookie() sessions.Store {
	store := cookie.NewStore([]byte("secret"))
	return store
}

func GetAllSession() []string {
	objects := reflect.TypeOf(Table{})
	for i := 0; i < objects.NumField(); i++ {
		fieldName := objects.Field(i).Name
		session_name_list = append(session_name_list, strings.ToLower(fieldName[:len(fieldName)-7]))
	}
	return session_name_list
}

func RegisterAll(objects ...interface{}) {
	for _, obj := range objects {
		gob.Register(obj)
	}
}
