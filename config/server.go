package config

import (
	"gopkg.in/yaml.v2"
	"log"
	"net/url"
	"os"
	"sync/atomic"
)

// ServerConfig represents the structure of the server configuration.
type ServerConfig struct {
	Servers []struct {
		URL string `yaml:"url"`
	} `yaml:"servers"`
}

// ServerPool holds the list of servers and the index for the next server.
type ServerPool struct {
	servers []*url.URL
	current uint64
}

// Global variable to hold the server pool.
var Pool ServerPool

// init reads the server configuration and initializes the server pool.
func init() {
	var conf ServerConfig
	data, err := os.ReadFile("config/server.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	for _, s := range conf.Servers {
		parsedURL, err := url.Parse(s.URL)
		if err != nil {
			log.Fatalf("error parsing URL: %v", err)
		}
		Pool.servers = append(Pool.servers, parsedURL)
	}
}

// NextIndex atomically increments the current index and returns the next server's index.
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, 1) % uint64(len(s.servers)))
}

// GetNextPeer returns the URL of the next server in the pool.
func (s *ServerPool) GetNextPeer() *url.URL {
	next := s.NextIndex()
	return s.servers[next]
}
