package app

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/repositories"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/httpserver"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func Run(cfg *config.Config) {
	if err := migrate(cfg.PG.URL); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	conn, err := pgxpool.New(context.Background(), cfg.PG.URL)
	if err != nil {
		log.Fatal("connection: ", err)
	}

	q := queries.New(conn)

	noteRepo := repositories.NewNoteRepo(q, conn)
	userRepo := repositories.NewUserRepo(q, conn)

	noteService := service.NewNoteService(noteRepo)
	userService := service.NewUserService(userRepo)

	httpserver := httpserver.New(cfg.HTTP.PORT)

	http.NewRoute(httpserver.App, noteService, userService)

	httpserver.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Shutdown")

	case err := <-httpserver.Notify():
		log.Panic("Httpserver: %v", err)
	}

	err = httpserver.Shutdown()
	if err != nil {
		log.Fatalf("Httpserver: %v", err)
	}
}

func migrate(url string) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	if err := goose.Up(db, "sql/migrations"); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
