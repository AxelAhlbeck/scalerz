package middlewares

import (
	"log"
	"math/rand"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const validToken = "superdupersecret"

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("%s %s", r.Method, r.RequestURI)
	})
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := authHeader[len(prefix):]
		if token != validToken {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RateLimitingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	limiter := time.Tick(100 * time.Millisecond)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-limiter
		next.ServeHTTP(w, r)
	})
}

func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := rand.Intn(100000)
		w.Header().Set("X-Request-ID", strconv.Itoa(requestID))
		next.ServeHTTP(w, r)
	})
}
