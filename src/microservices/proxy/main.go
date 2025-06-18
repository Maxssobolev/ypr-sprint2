package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

var (
	moviesMigrationPercent int
	monolithURL            *url.URL
	moviesServiceURL       *url.URL
)

func main() {
	// Инициализация переменных окружения
	initEnv()

	// Настройка reverse proxy
	monolithProxy := httputil.NewSingleHostReverseProxy(monolithURL)
	moviesProxy := httputil.NewSingleHostReverseProxy(moviesServiceURL)

	// Обработчики маршрутов
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/movies", moviesHandler(monolithProxy, moviesProxy))
	http.Handle("/", monolithProxy) // Все остальные запросы направляем в монолит

	port := getEnv("PORT", "8082")
	log.Printf("Starting proxy server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initEnv() {
	// Получение процента миграции
	migrationPercent, _ := strconv.Atoi(getEnv("MOVIES_MIGRATION_PERCENT", "0"))
	if migrationPercent < 0 || migrationPercent > 100 {
		log.Fatal("MOVIES_MIGRATION_PERCENT must be between 0 and 100")
	}

	// URL монолита
	monolith, _ := url.Parse(getEnv("MONOLITH_URL", "http://localhost:8080"))
	monolithURL = monolith

	// URL сервиса фильмов
	moviesSvc, _ := url.Parse(getEnv("MOVIES_SERVICE_URL", "http://localhost:8081"))
	moviesServiceURL = moviesSvc

	log.Printf("Migration config: %d%% traffic to movies service", migrationPercent)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func moviesHandler(monolithProxy, moviesProxy *httputil.ReverseProxy) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Определение целевого сервиса на основе feature flag
		target := monolithProxy
		if shouldRouteToMoviesService() {
			log.Printf("Routing to movies service: %s %s", r.Method, r.URL.Path)
			target = moviesProxy
		} else {
			log.Printf("Routing to monolith: %s %s", r.Method, r.URL.Path)
		}

		// Копирование оригинальных заголовков
		r.Header.Set("X-Forwarded-Host", r.Host)
		r.Header.Set("X-Proxy-Service", "cinema-proxy")

		// Перенаправление запроса
		target.ServeHTTP(w, r)
	}
}

func shouldRouteToMoviesService() bool {
	return rand.Intn(100) < moviesMigrationPercent
}
