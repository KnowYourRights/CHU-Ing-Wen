package main

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	rateInterval           = 69
	tokenBursts            = 42069
	freeRealEstateInterval = 42069
	rateLimitExceedErrStr  = "（４２９）泥讓婐休息１下毫ㄇ？？？"
)

var visitorAddrs = make(map[string]*rate.Limiter)
var mu sync.Mutex

func init() {
	go func() {
		for {
			time.Sleep(freeRealEstateInterval * time.Second)

			mu.Lock()
			visitorAddrs = make(map[string]*rate.Limiter)
			mu.Unlock()
		}
	}()
}

func getLimiter(addr string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	lim, ok := visitorAddrs[addr]
	if !ok {
		newLim := rate.NewLimiter(rateInterval, tokenBursts)
		visitorAddrs[addr] = newLim
		return newLim
	}
	return lim
}

func rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the IP address of the request.
		addr := r.Header.Get("X-FORWARDED-FOR")
		if addr == "" {
			addr = r.RemoteAddr
		}

		limiter := getLimiter(addr)
		if !limiter.Allow() {
			http.Error(w, rateLimitExceedErrStr, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
