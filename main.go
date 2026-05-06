package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"dong/internal/config"
	"dong/internal/router"
)

//go:embed web
var embeddedFiles embed.FS

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 静态目录
	webFS, err := fs.Sub(embeddedFiles, "web")
	if err != nil {
		log.Fatalf("failed to sub embed fs: %v", err)
	}

	r, err := router.Setup(cfg, http.FS(webFS))
	if err != nil {
		log.Fatalf("failed to setup router: %v", err)
	}

	log.Printf("SERVER_PORT: %s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
