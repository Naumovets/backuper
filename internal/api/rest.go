package api

import (
	"net/http"

	"github.com/Naumovets/backuper/internal/services"
	"github.com/gin-gonic/gin"
)

func InitRest(servicer services.Servicer) *http.Server {
	r := gin.Default()
	r.Use(traceLogMiddleware())

	h := NewHandler(servicer)

	r.GET("/health", h.Health)
	r.GET("/postgres/backupers", h.GetPostgresBackupers)
	r.GET("/logging/config", h.GetLoggingConfig)
	r.GET("/storage/config", h.GetStorageConfig)

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}
