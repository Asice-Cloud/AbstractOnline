package session

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func Test_init(t *testing.T) {
	InitSession(gin.New())
}
