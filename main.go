package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/misshanya/secret-santa/db"
	"github.com/misshanya/secret-santa/routes/auth"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	serverPort := os.Getenv("SERVER_PORT")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("could not connect to db")
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	authApi := auth.NewAuthAPI(queries)

	r.Post("/register", authApi.RegisterUser)

	fmt.Println("Server is up")

	http.ListenAndServe(fmt.Sprintf(":%v", serverPort), r)
}
