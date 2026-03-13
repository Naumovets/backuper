package api

import (
	"net/http"

	"github.com/Naumovets/backuper/internal/logger"
	"github.com/Naumovets/backuper/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type handler struct {
	servicer services.Servicer
}

func NewHandler(servicer services.Servicer) *handler {
	return &handler{
		servicer: servicer,
	}
}

func (r *handler) Health(c *gin.Context) {
	c.String(http.StatusOK, "alive")
}

func (r *handler) GetPostgresBackupers(c *gin.Context) {
	log := logger.GetLogger(c.Request.Context())
	backuper, err := r.servicer.GetPostgresBackupers(c.Request.Context())
	if err != nil {
		log.Error("cannot get postgres backupers", zap.Error(err))

		c.Status(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, backuper)
}

func (r *handler) GetLoggingConfig(c *gin.Context) {
	log := logger.GetLogger(c.Request.Context())
	backuper, err := r.servicer.GetLoggingConfig(c.Request.Context())
	if err != nil {
		log.Error("cannot get logging config", zap.Error(err))

		c.Status(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, backuper)
}

func (r *handler) GetStorageConfig(c *gin.Context) {
	log := logger.GetLogger(c.Request.Context())
	backuper, err := r.servicer.GetStorageConfig(c.Request.Context())
	if err != nil {
		log.Error("cannot get storage config", zap.Error(err))

		c.Status(http.StatusBadRequest)
	}

	c.JSON(http.StatusOK, backuper)
}
