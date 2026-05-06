package handler

import (
	"dong/internal/config"

	"github.com/gin-gonic/gin"
)

func NewHealthHandler(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"dbHost": cfg.DBHost,
		})
	}
}
