package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

type us struct {
	UserID       int
	UserName     string
	AccessToken  string
	RefreshToken string
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	store := memstore.NewStore([]byte("secret"))
	names := []string{"user", "admin"}
	r.Use(sessions.SessionsMany(names, store))
	return r
}

func TestOperations(t *testing.T) {
	router := setupRouter()

	// Set up a sample route to test session operations
	router.GET("/test", func(c *gin.Context) {
		SessionSet("user", c, "user", us{UserID: 1, UserName: "testuser", AccessToken: "token", RefreshToken: "refresh"})
		user := SessionGet("user", c, "user")
		SessionUpdate("user", c, "user", us{UserID: 1, UserName: "updateduser", AccessToken: "newtoken", RefreshToken: "newrefresh"})
		SessionDelete("user", c, "user")
		c.JSON(200, user)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)

	// Send request to the router
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}

	// Check response body
	expectedUser := `{"UserID":1,"UserName":"testuser","AccessToken":"token","RefreshToken":"refresh"}`
	if w.Body.String() != expectedUser {
		t.Fatalf("Expected body %s, but got %s", expectedUser, w.Body.String())
	}
}
