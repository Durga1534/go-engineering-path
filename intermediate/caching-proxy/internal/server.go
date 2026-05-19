package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"caching-proxy/cache"
)

type Server struct {
	Origin string
	Cache  *cache.Store
}

func NewServer(origin string, store *cache.Store) *Server {
	return &Server{
		Origin: strings.TrimSuffix(origin, "/"),
		Cache:  store,
	}
}

func (p *Server) HandleProxy(w http.ResponseWriter, r *http.Request) {
	cacheKey := fmt.Sprintf("%s:%s", r.Method, r.URL.RequestURI())

	if cachedRes, hit := p.Cache.Get(cacheKey); hit {
		log.Printf("[CACHE HIT] Serving path: %s", r.URL.RequestURI())

		for key, values := range cachedRes.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(cachedRes.StatusCode)
		w.Write(cachedRes.Body)
		return
	}

	log.Printf("[CACHE MISS] Forwarding path to origin: %s", r.URL.RequestURI())
	originURL := fmt.Sprintf("%s%s", p.Origin, r.URL.RequestURI())

	var bodyBytes []byte
	if r.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read incoming request body", http.StatusInternalServerError)
			return
		}
	}
	originReq, err := http.NewRequest(r.Method, originURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to create origin request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			originReq.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(originReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Origin server unreachable: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read origin response content", http.StatusInternalServerError)
		return
	}
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	if resp.StatusCode == http.StatusOK {
		p.Cache.Set(cacheKey, cache.CachedResponse{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       respBody,
		})
	}
}
