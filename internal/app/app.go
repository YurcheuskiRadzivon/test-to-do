package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/admin"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/auth"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/file"
	middleware "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/middleware/auth"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/note"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/user"
	authmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/auth"
	encryptmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/encrypt"
	filemanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/file"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/repositories"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/migrations"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/generator"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/httpserver"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/jackc/pgx/v5/pgxpool"
	migrator "github.com/rubenv/sql-migrate"

	_ "github.com/lib/pq"
)

func Run(cfg *config.Config) {
	//PG
	migrator := migrations.NewPGMigrator(cfg.PG.URL, "sql/migrations")
	if err := migrator.Migrate(); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	conn, err := pgxpool.New(context.Background(), cfg.PG.URL)
	if err != nil {
		log.Fatal("connection: ", err)
	}

	q := queries.New(conn)

	//Generator
	g := generator.NewGenerator()

	//JWT
	jwtS := jwtservice.New(cfg.JWT.SECRETKEY)

	//Managers
	authManager := authmanage.NewAuthManage(jwtS)
	encryptManager := encryptmanage.NewEncrypter()
	fileManager := filemanage.NewFileManage(g)

	//Repo
	noteRepo := repositories.NewNoteRepo(q, conn)
	userRepo := repositories.NewUserRepo(q, conn)
	fileMetaRepo := repositories.NewFileMetaRepo(q, conn)

	//Service
	noteService := service.NewNoteService(noteRepo)
	userService := service.NewUserService(userRepo)
	fileMetaService := service.NewFileMetaService(fileMetaRepo)

	//Middleware
	authMiddleware := middleware.NewAuthMW(fileMetaService, authManager, userService, cfg)

	//Controller
	authController := auth.NewAuthControl(userService, authManager, encryptManager)
	userController := user.NewUserControl(userService, authManager, encryptManager)
	adminController := admin.NewAdminControl(userService, authManager, encryptManager)
	noteController := note.NewNoteControl(fileMetaService, noteService, authManager, fileManager)
	fileController := file.NewFileControl(fileMetaService, fileManager, authManager, noteService)

	httpserver := httpserver.New(cfg.HTTP.PORT)

	http.NewRoute(
		httpserver.App,
		noteController,
		userController,
		adminController,
		authController,
		authMiddleware,
		fileController,
		cfg,
	)

	httpserver.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case <-interrupt:
		log.Println("Shutdown")

	case err := <-httpserver.Notify():
		log.Panicf("Httpserver: %s", err)
	}

	err = httpserver.Shutdown()
	if err != nil {
		log.Fatalf("Httpserver: %v", err)
	}
}

func migrate(url string) error {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer db.Close()

	migrations := &migrator.FileMigrationSource{
		Dir: "sql/migrations",
	}

	n, err := migrator.Exec(db, "postgres", migrations, migrator.Up)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}
