package http

import (
	"net/http"

	"github.com/dilroop-us/ecommerce-go/internal/platform/requestid"
	"github.com/dilroop-us/ecommerce-go/internal/product"
	"github.com/gin-gonic/gin"
)

func Router(store *product.Store) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestid.Middleware())

	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })

	r.GET("/products", func(c *gin.Context) {
		c.JSON(http.StatusOK, store.List())
	})

	type createReq struct {
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price" binding:"required,gt=0"`
	}
	r.POST("/products", func(c *gin.Context) {
		var req createReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := store.Create(req.Name, req.Price)
		c.JSON(http.StatusCreated, p)
	})

	return r
}
