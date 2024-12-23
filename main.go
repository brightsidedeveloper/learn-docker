package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Woohoo, Kubernetes with CI!!!"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Fire!"))
	})

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {

		product := product{
			ID:    1,
			Name:  "Laptop",
			Price: 1000,
		}

		if err := json.NewEncoder(w).Encode(product); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env variable is not set")
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: withCORS(mux),
	}

	go func() {
		log.Printf("Starting server on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("Gracefully shutting down server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server shutdown successfully")
}
