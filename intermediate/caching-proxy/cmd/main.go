package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"caching-proxy/cache"
	"caching-proxy/internal"
)

func main() {
	port := flag.Int("port", 0, "The port on which proxy server will run")
	origin := flag.String("origin", "", "The URL of the server to which requests will be forwarded")
	clearCache := flag.Bool("clear-cache", false, "Clear  the proxy cache and exit")

	flag.Parse()

	if *clearCache {
		fmt.Println("Cache cleared.")
		os.Exit(0)
	}
	if *port == 0 || *origin == "" {
		fmt.Println("Error: Missing required arguments")
		fmt.Println("Usage: caching-proxy --port <number> --origin<url>")
		os.Exit(1)
	}
	sharedCacheStore := cache.NewStore()
	proxyEngine := internal.NewServer(*origin, sharedCacheStore)

	mux := http.NewServeMux()
	mux.HandleFunc("/", proxyEngine.HandleProxy)

	serverAddr := fmt.Sprintf(":%d", *port)
	log.Printf("Starting caching proxy server on port %d...", *port)
	log.Printf("Forwarding traffic to origin: %s", *origin)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server lifecycle failed to start up :%v", err)
	}
}
