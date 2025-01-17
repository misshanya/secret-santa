package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/misshanya/secret-santa/config"
	"github.com/misshanya/secret-santa/db"
	"github.com/misshanya/secret-santa/middlewares"
	"github.com/misshanya/secret-santa/routes"
)

func main() {
	cfg := config.GetConfig()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("could not connect to db")
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	authApi := routes.NewAuthAPI(queries)
	roomsApi := routes.NewRoomsAPI(queries)
	participantsAPI := routes.NewParticipantsAPI(queries, conn)

	r.Post("/register", authApi.RegisterUser)
	r.Post("/login", authApi.Login)
	r.With(middlewares.Auth).Get("/me", authApi.GetMyInfo)

	r.With(middlewares.Auth).Post("/rooms", roomsApi.CreateRoom)
	r.With(middlewares.Auth).Get("/rooms", roomsApi.MyRooms)
	r.With(middlewares.Auth).Delete("/rooms/{id}", roomsApi.DeleteRoom)

	r.With(middlewares.Auth).Post("/rooms/{id}/participants", participantsAPI.NewParticipant)
	r.With(middlewares.Auth).Delete("/rooms/{id}/participants/me", participantsAPI.DeleteParticipant)
	r.With(middlewares.Auth).Get("/rooms/{id}/wish", participantsAPI.GetWish)
	r.With(middlewares.Auth).Patch("/rooms/{id}/wish", participantsAPI.SetWish)

	r.With(middlewares.Auth).Get("/rooms/{id}/participants", participantsAPI.GetParticipants)
	r.With(middlewares.Auth).Post("/rooms/{id}/participants/distribute", participantsAPI.DistributeParticipants)

	fmt.Println("Server is up")

	http.ListenAndServe(fmt.Sprintf(":%v", cfg.ServerPort), r)
}
