package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/redis"
)

func inRedis() sessions.Store {
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	return store
}

func inCookie() sessions.Store {
	store := cookie.NewStore([]byte("secret"))
	return store
}
