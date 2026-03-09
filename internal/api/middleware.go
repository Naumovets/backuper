package api

import (
	"context"

	"github.com/Naumovets/Backuper/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func traceLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		trace_id, ok := ctx.Value("trace_id").(uuid.UUID)
		if !ok {
			trace_id = uuid.New()
			ctx = context.WithValue(ctx, "trace_id", trace_id)
		}

		log := zap.L().With(zap.String("trace_id", trace_id.String()))

		ctx = context.WithValue(ctx, logger.CtxKeyLogger, log)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
