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
       if fileServer == nil {
              return nil, fmt.Errorf("router setup failed: fileServer is nil")
       }

       // 注册静态首页
       r.GET("/", gin.WrapH(fileServer))
       // 注册健康检查
       r.GET("/api/health", handler.NewHealthHandler(cfg))
       r.GET("/api/mysql/testConnect", handler.TestConnect)
       r.GET("/api/mysql/version", handler.Version)
       // 注册兜底静态
       r.NoRoute(gin.WrapH(fileServer))

       return r, nil
}
