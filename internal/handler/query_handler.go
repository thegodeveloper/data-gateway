package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryRequest struct {
	Source string                 `json:"source"`
	Query  map[string]interface{} `json:"query"`
}

func HandleQuery(sources map[string]DataSource) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req QueryRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		ds, ok := sources[req.Source]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown data source"})
			return
		}

		result, err := ds.Query(c.Request.Context(), req.Query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
