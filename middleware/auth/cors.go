package auth

import (
	"github.com/gin-contrib/cors"
)

func CorsInit() cors.Config {
	// set up CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}                                        // allowed origin，use * represent for plural
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // allowed http method
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // allowed http header
	return config
}
