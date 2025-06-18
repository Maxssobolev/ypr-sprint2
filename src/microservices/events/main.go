package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"events/api"
	"events/config"
	"events/kafka"
)

func main() {
	cfg := config.Load()

	// Инициализация продюсера
	producer := kafka.NewProducer(cfg.KafkaBrokers)
	defer producer.Close()

	// Инициализация и запуск consumer
	consumer := kafka.NewConsumer([]string{cfg.KafkaBrokers}, "events-service-group")
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go consumer.Consume(ctx, &wg)

	router := api.SetupRouter(producer)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

	<-done
	log.Println("Server stopped")

	cancel()
	wg.Wait()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server shutdown failed:", err)
	}
	log.Println("Server exited properly")
}
