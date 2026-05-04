package router

import (
	"fmt"
	"net/http"

	"dong/internal/config"
	"dong/internal/handler"
	"github.com/gin-gonic/gin"
)

func Setup(cfg config.Config, staticFS http.FileSystem) (*gin.Engine, error) {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	if err := r.SetTrustedProxies(cfg.TrustedProxies); err != nil {
		return nil, fmt.Errorf("set trusted proxies: %w", err)
	}

	fileServer := http.FileServer(staticFS)
	r.GET("/", gin.WrapH(fileServer))
	r.GET("/api/health", handler.NewHealthHandler(cfg))
	r.NoRoute(gin.WrapH(fileServer))

	return r, nil
}
