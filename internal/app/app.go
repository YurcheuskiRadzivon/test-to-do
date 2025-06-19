package app

import (
	"context"
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
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/storages"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/database/queries"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/infrastructure/migrations"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/generator"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/httpserver"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go"

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

	//MinioClient
	//Create minio client
	minioClient, err := minio.New(
		cfg.LOCALSTACK.INTERNAL_ENDPOINT,
		cfg.LOCALSTACK.ACCESS_KEY,
		cfg.LOCALSTACK.SECRET_KEY,
		false,
	)
	if err != nil {
		log.Fatal("minio: ", err)
	}

	//Create bucket
	exists, err := minioClient.BucketExists(cfg.LOCALSTACK.BUCKET)
	if err != nil {
		log.Fatal("minio: create bucket: ", err)
	}
	if !exists {
		err := minioClient.MakeBucket(cfg.LOCALSTACK.BUCKET, "")
		if err != nil {
			log.Fatal("minio: create bucket: ", err)
		}
	}

	//Debug connection test
	buckets, err := minioClient.ListBuckets()
	if err != nil {
		log.Fatalf("minio: cannot connect to s3: %v", err)
	}
	log.Println("Succesfully connections. Buckets:", buckets)

	//Repo
	noteRepo := repositories.NewNoteRepo(q, conn)
	userRepo := repositories.NewUserRepo(q, conn)
	fileMetaRepo := repositories.NewFileMetaRepo(q, conn)

	//Storage
	s3Storage := storages.NewS3Storage(
		minioClient,
		cfg.LOCALSTACK.BUCKET,
		cfg.LOCALSTACK.EXTERNAL_ENDPOINT,
		cfg.LOCALSTACK.INTERNAL_ENDPOINT,
	)

	_ = s3Storage
	fsStorage := storages.NewFSStorage(
		cfg.FSSTORAGE.PATH,
		cfg.FSSTORAGE.EXTERNAL_ENDPOINT,
		cfg.APP.DOMAIN,
	)

	//Managers
	authManager := authmanage.NewAuthManage(jwtS)
	encryptManager := encryptmanage.NewEncrypter()
	fileManager := filemanage.NewFileManage(g, fsStorage)

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
