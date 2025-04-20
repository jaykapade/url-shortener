package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jaykapade/url-shortener/internal/auth"
	"github.com/jaykapade/url-shortener/internal/db"
	"github.com/jaykapade/url-shortener/internal/shortener"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file:", err)
		return
	}
	log.Println("Loaded .env file")

	DB_URL := os.Getenv("DB_URL")
	conn, err := pgx.Connect(context.Background(), DB_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	queries := db.New(conn)

	authHandler := &auth.AuthHandler{
		DB: queries,
	}
	shortenerHandler := &shortener.ShortenerHandler{
		DB: queries,
	}

	// Initialize Sonyflake Generator
	shortener.InitIDGenerator()

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		// Auth routes
		r.Post("/login", authHandler.LoginHandler)
		r.Post("/register", authHandler.RegisterHandler)
		r.With(auth.JWTMiddleware).Get("/test", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value(auth.UserIDKey).(string)
			fmt.Println("User ID:", userID)
			w.WriteHeader(http.StatusOK)
		})
		// Shortener routes
		r.With(auth.JWTMiddleware).Post("/shortener", shortenerHandler.CreateShortCodeHandler)
		r.Get("/{short_code}", shortenerHandler.RedirectLinkHandler)
	})

	log.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))

}
