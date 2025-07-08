package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"lkgiovani_go_rinha_2025/internal"
	"lkgiovani_go_rinha_2025/internal/config"

	"github.com/panjf2000/gnet/v2"
)

func main() {
	var configPath string
	var port int
	var multicore bool

	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.IntVar(&port, "port", 8080, "Server port")
	flag.BoolVar(&multicore, "multicore", true, "Enable multicore")
	flag.Parse()

	// Load config if provided
	var cfg config.Config
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			cfg.LoadConfig(configPath)
			port = cfg.AppPort
		}
	}

	addr := fmt.Sprintf("tcp://0.0.0.0:%d", port)

	server := &internal.HttpServer{
		Addr:      addr,
		Multicore: multicore,
	}

	log.Printf("Starting server on %s (multicore: %t)", addr, multicore)
	log.Println("Server exits:", gnet.Run(server, addr, gnet.WithMulticore(multicore)))
}
