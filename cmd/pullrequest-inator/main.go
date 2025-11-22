package main

import (
	"context"
	"log"
	"os"
	"pullrequest-inator/internal/api"
	pg2 "pullrequest-inator/internal/infrastructure/repositories/pg"
	"pullrequest-inator/internal/infrastructure/services"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	port = strings.TrimSpace(port)
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	log.Printf("Starting server on %s (raw value was: %q)", port, os.Getenv("SERVER_PORT"))
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Failed to connect to database:", err)
	}
	defer pool.Close()

	prRepo := pg2.NewPullRequestRepository(pool)
	statusRepo := pg2.NewStatusRepository(pool)
	teamRepo := pg2.NewTeamRepository(pool)
	userRepo := pg2.NewUserRepository(pool)

	prService, err := services.NewPullRequestService(userRepo, prRepo, teamRepo, statusRepo)
	if err != nil {
		log.Fatal(err)
	}
	teamService, err := services.NewTeamService(teamRepo, userRepo)
	if err != nil {
		log.Fatal(err)
	}
	userService, err := services.NewUserService(userRepo)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	server, err := api.NewServer(prService, teamService, userService)
	if err != nil {
		log.Fatal(err)
	}

	api.RegisterHandlers(e, server)

	if err := e.Start(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}
