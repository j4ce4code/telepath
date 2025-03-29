// telepath.go
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"sync"
	"syscall"
)

type Config struct {
	Mode       string            `json:"mode"`       // "header" or "path"
	HeaderName string            `json:"headerName"` // e.g., "X-Runtime-Env"
	Routes     map[string]string `json:"routes"`
}

type ProxyServer struct {
	mu     sync.RWMutex
	config *Config
	path   string
}

func (ps *ProxyServer) LoadConfig() error {
	file, err := os.Open(ps.path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, _ := io.ReadAll(file)
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	ps.mu.Lock()
	ps.config = &cfg
	ps.mu.Unlock()
	return nil
}

func (ps *ProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ps.mu.RLock()
	cfg := ps.config
	ps.mu.RUnlock()

	var key string
	if cfg.Mode == "header" {
		key = r.Header.Get(cfg.HeaderName)
		if key == "" {
			http.Error(w, "Missing header", http.StatusBadRequest)
			return
		}
	} else if cfg.Mode == "path" {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 1 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}
		key = parts[0]
		r.URL.Path = "/" + path.Join(parts[1:]...) // Strip first segment
	}

	target, ok := cfg.Routes[key]
	if !ok {
		http.Error(w, "Unknown route key", http.StatusNotFound)
		return
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		http.Error(w, "Bad target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	cfgPath := "./telepath.json"
	server := &ProxyServer{path: cfgPath}

	if err := server.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Listen for SIGHUP
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for range c {
			log.Println("Reloading config...")
			if err := server.LoadConfig(); err != nil {
				log.Printf("Error reloading config: %v", err)
			} else {
				log.Println("Config reloaded.")
			}
		}
	}()

	log.Println("telepath proxy starting on :8080...")
	http.ListenAndServe(":8080", server)
}
