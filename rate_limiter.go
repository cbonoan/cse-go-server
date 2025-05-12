package main

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	mu sync.Mutex
	visitors = make(map[string]*rate.Limiter)
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]

	if !exists {
		// Rate limit: 1 request per 5 seconds with a burst of 3
		limiter = rate.NewLimiter(rate.Every(5*time.Second), 3)
		visitors[ip] = limiter

		// Remove IP from map after 1 minute
		go func() {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			delete(visitors, ip)
			mu.Unlock()
		}()
	}

	return limiter
}

func RateLimiter(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
		limiter := getVisitor(ip)

        if !limiter.Allow() {
            message := Response{
                Message:   "The API is at capacity, try again later.",
                ResponseCode: http.StatusTooManyRequests,
            }

            w.WriteHeader(http.StatusTooManyRequests)
            json.NewEncoder(w).Encode(&message)
            return
        }
        next.ServeHTTP(w, r)
    })
}
