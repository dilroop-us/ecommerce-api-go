package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dilroop-us/ecommerce-go/internal/db"
	"github.com/dilroop-us/ecommerce-go/internal/product"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql
)

func main() {
	// Read DB URL from environment variable
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	// Open DB connection using pgx stdlib
	sqldb, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	defer sqldb.Close()

	// Verify DB connection
	if err := sqldb.Ping(); err != nil {
		log.Fatalf("db ping: %v", err)
	}
	log.Println("✅ Connected to Postgres")

	// Initialize generated queries + layers
	q := db.New(sqldb)
	repo := product.NewRepository(q)
	svc := product.NewService(repo)

	// Setup Gin router
	r := gin.New()
	r.Use(gin.Recovery())

	// Basic health route
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// List all products
	r.GET("/products", func(c *gin.Context) {
		items, err := svc.List(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// Create new product
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

		item, err := svc.Create(c.Request.Context(), req.Name, req.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	// HTTP server setup
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start server in background
	go func() {
		log.Println("🚀 Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}
