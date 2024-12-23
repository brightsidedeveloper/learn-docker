package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

/*

 - name: Log in to Docker Hub
        run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

      - name: Build and Push Docker Image
        run: |
          docker build -t ${{ secrets.DOCKER_USERNAME }}/go-server:latest .
          docker push ${{ secrets.DOCKER_USERNAME }}/go-server:latest


*/

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Tim!!"))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Fire!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env variable is not set")
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
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
