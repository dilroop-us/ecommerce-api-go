package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	apphttp "github.com/dilroop-us/ecommerce-go/internal/http"
	"github.com/dilroop-us/ecommerce-go/internal/platform/logger"
	"github.com/dilroop-us/ecommerce-go/internal/product"
)

func main() {
	zl, err := logger.New()
	if err != nil {
		log.Fatalf("logger: %v", err)
	}
	defer zl.Sync()

	store := product.NewStore()
	handler := apphttp.Router(store)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		zl.Info("server starting", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zl.Fatal("listen", zap.Error(err))
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zl.Error("server shutdown", zap.Error(err))
	}
	zl.Info("server exited")
}
