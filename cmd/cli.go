// cmd/cli.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CLIConfig struct {
	Mode       string            `json:"mode"`
	HeaderName string            `json:"headerName"`
	Routes     map[string]string `json:"routes"`
}

const configPath = "./telepath.json"

func loadCLIConfig() (*CLIConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg CLIConfig
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func writeCLIConfig(cfg *CLIConfig) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func listRoutes() {
	cfg, err := loadCLIConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("Current routes:")
	for k, v := range cfg.Routes {
		fmt.Printf("  %s -> %s\n", k, v)
	}
}

func addRoute(key, target string) {
	cfg, err := loadCLIConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if cfg.Routes == nil {
		cfg.Routes = map[string]string{}
	}
	cfg.Routes[key] = target
	if err := writeCLIConfig(cfg); err != nil {
		log.Fatalf("Failed to write config: %v", err)
	}
	fmt.Printf("Added route: %s -> %s\n", key, target)
}

func removeRoute(key string) {
	cfg, err := loadCLIConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	delete(cfg.Routes, key)
	if err := writeCLIConfig(cfg); err != nil {
		log.Fatalf("Failed to write config: %v", err)
	}
	fmt.Printf("Removed route: %s\n", key)
}

func refreshServer() {
	cmd := exec.Command("pkill", "-SIGHUP", "telepath")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to send SIGHUP: %v", err)
	}
	fmt.Println("Sent SIGHUP to telepath server")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: telepath <command> [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "route":
		if len(os.Args) < 3 {
			log.Fatal("Usage: telepath route <add|remove|list> ...")
		}
		sub := os.Args[2]
		switch sub {
		case "add":
			if len(os.Args) != 5 {
				log.Fatal("Usage: telepath route add <key> <targetURL>")
			}
			addRoute(os.Args[3], os.Args[4])
		case "remove":
			if len(os.Args) != 4 {
				log.Fatal("Usage: telepath route remove <key>")
			}
			removeRoute(os.Args[3])
		case "list":
			listRoutes()
		default:
			log.Fatalf("Unknown subcommand: %s", sub)
		}
	case "refresh":
		refreshServer()
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}
