package balance

import (
	"Abstract/config"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

var serverPool config.ServerPool

// NextIndex 循环获取下一个服务器索引
// LoadBalancerMiddleware creates a middleware for load balancing

func LoadBalancerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		peer := config.Pool.GetNextPeer()                 // Get the next server URL from the pool
		proxy := httputil.NewSingleHostReverseProxy(peer) // Create a reverse proxy to that server

		// Update the request to reflect the scheme and host of the selected peer
		originalHost := c.Request.Host
		c.Request.URL.Scheme = peer.Scheme
		c.Request.URL.Host = peer.Host
		c.Request.Header.Set("X-Forwarded-Host", originalHost)
		c.Request.Host = peer.Host

		// Serve the request via the reverse proxy
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
