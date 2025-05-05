// Package http
// internal/transport/http/server.go
package http

import (
	"github.com/thegodeveloper/data-gateway/internal/app"
	"github.com/thegodeveloper/data-gateway/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func StartServer(svc *app.GatewayService, port string) {
	r := gin.Default()
	r.Use(otelgin.Middleware("data-gateway"))

	r.POST("/query", func(c *gin.Context) {
		var req domain.QueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		res, err := svc.HandleQuery(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.Run(":" + port)
}
