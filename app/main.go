package main

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/namishh/holmes/database"
	"github.com/namishh/holmes/handlers"
	"github.com/namishh/holmes/services"
)

func initMinioClient() (*minio.Client, error) {
	endpoint := os.Getenv("BUCKET_ENDPOINT")
	accessKeyID := os.Getenv("BUCKET_ACCESSKEY")
	secretAccessKey := os.Getenv("BUCKET_SECRETKEY")
	useSSL := true
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to the bucket")

	return minioClient, nil
}

func main() {
	if os.Getenv("ENVIRONMENT") == "DEV" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	minioClient, err := initMinioClient()
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	e := echo.New()
	SECRET_KEY := os.Getenv("SECRET")
	DB_NAME := os.Getenv("DB_NAME")
	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SECRET_KEY))))

	e.Static("/static", "public")

	store, err := database.NewDatabaseStore(DB_NAME)
	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	us := services.NewUserService(services.User{}, store, minioClient)
	ah := handlers.NewAuthHandler(us)

	handlers.SetupRoutes(e, ah)

	// Start server
	e.Logger.Fatal(e.Start(":4200"))
}
