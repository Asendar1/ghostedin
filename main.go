package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	InitDB()

	// Home
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "home/", http.StatusPermanentRedirect)
	})
	r.Get("/home/", GetHomePage)

	r.Route("/track", func(r chi.Router) {
		r.Get("/", LoadSheet)
		r.Get("/applications", GetRows)
		r.Get("/statistics", GetStatistics)
	})



	//#region Server start & Graceful Shutdown

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func ()  {
		log.Fatal(srv.ListenAndServe())
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<- quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Println("Server stopped")
	//#endregion
}
