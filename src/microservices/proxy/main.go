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
	"time"
	"strings"
)

var (
	migrationPercent int
	monolithURL            *url.URL
	moviesServiceURL       *url.URL
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// Инициализация переменных окружения
	initEnv()

	// Настройка reverse proxy
	monolithProxy := httputil.NewSingleHostReverseProxy(monolithURL)
	moviesProxy := httputil.NewSingleHostReverseProxy(moviesServiceURL)

	// Обработчики маршрутов
	http.HandleFunc("/health", healthHandler)

    http.Handle("/", proxyHandler(monolithProxy, moviesProxy))

	port := getEnv("PORT", "8082")
	log.Printf("Starting proxy server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func initEnv() {
	// Получение процента миграции
	migrationPercent, _ = strconv.Atoi(getEnv("MOVIES_MIGRATION_PERCENT", "0"))
	if migrationPercent < 0 || migrationPercent > 100 {
		log.Fatal("MOVIES_MIGRATION_PERCENT must be between 0 and 100")
	}

	// URL монолита
	monolith, _ := url.Parse(getEnv("MONOLITH_URL", "http://localhost:8010"))
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

func proxyHandler(monolithProxy, moviesProxy *httputil.ReverseProxy) http.HandlerFunc {
	 return func(w http.ResponseWriter, r *http.Request) {
        var target *httputil.ReverseProxy

        path := strings.TrimSuffix(r.URL.Path, "/") // для /api/movies/ → /api/movies

        // Пути, которые всегда идут на монолит
        if path == "/api/users" ||
            path == "/api/payments" ||
            path == "/api/subscriptions" {
            target = monolithProxy
        } else if path == "/api/movies" || path == "/api/movies/health" {
            // Маршруты, связанные с фильмами (учитывают процент миграции)
            if shouldRouteToMoviesService() {
                log.Printf("Routing to movies service: %s %s", r.Method, r.URL.Path)
                target = moviesProxy
            } else {
                log.Printf("Routing to monolith: %s %s", r.Method, r.URL.Path)
                target = monolithProxy
            }
        } else if path == "/health" {
            // Health check самого прокси
            healthHandler(w, r)
            return
        } else {
            // Все остальные запросы направляем в монолит
            log.Printf("Routing to monolith by default: %s %s", r.Method, r.URL.Path)
            target = monolithProxy
        }

        // Устанавливаем заголовки
        r.Header.Set("X-Forwarded-Host", r.Host)
        r.Header.Set("X-Proxy-Service", "cinema-proxy")

        // Проксируем запрос
        target.ServeHTTP(w, r)
    }
}

func shouldRouteToMoviesService() bool {
	return rand.Intn(100) < migrationPercent
}
