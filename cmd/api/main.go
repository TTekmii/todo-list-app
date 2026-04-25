package main

// @title           Todo List API
// @version         0.1.0
// @description     This is a sample Todo List server using Gin and Clean Architecture.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/TTekmii/todo-list-app/internal/app/auth"
	"github.com/TTekmii/todo-list-app/internal/app/todo"
	"github.com/TTekmii/todo-list-app/internal/infrastructure/database"
	"github.com/TTekmii/todo-list-app/internal/infrastructure/repository"
	"github.com/TTekmii/todo-list-app/internal/lib/logger"
	"github.com/TTekmii/todo-list-app/internal/lib/logger/sl"
	transportServer "github.com/TTekmii/todo-list-app/internal/transport/http-server"
	"github.com/TTekmii/todo-list-app/internal/transport/http-server/handler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("No .env file found, using system environment variables")
	}

	if err := initConfig(); err != nil {
		slog.Error("failed to initialize config", sl.Err(err))
	}

	cfgLogger := logger.Config{
		Level:  viper.GetString("logger.Level"),
		Format: viper.GetString("logger.format"),
		Env:    viper.GetString("app.env"),
	}
	log := logger.New(cfgLogger)

	log.Info("Starting application",
		slog.String("env", cfgLogger.Env),
		slog.String("log_level", cfgLogger.Level),
	)

	dbConfig := database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	log.Info("Database connected successfully")

	authRepo := repository.NewAuthPostgres(db)
	listRepo := repository.NewTodoListPostgres(db)
	itemRepo := repository.NewTodoItemPostgres(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = viper.GetString("jwt.secret")
	}
	if jwtSecret == "" {
		log.Error("JWT secret is not configured")
		os.Exit(1)
	}

	jwtTTL := viper.GetDuration("jwt.ttl")
	if jwtTTL == 0 {
		jwtTTL = 12 * time.Hour
	}

	bcryptCost := viper.GetInt("bcrypt.cost")
	if bcryptCost == 0 {
		bcryptCost = 10
	}

	log.Info("JWT config loaded",
		slog.Duration("ttl", jwtTTL),
		slog.Int("bcrypt_cost", bcryptCost),
	)

	authService := auth.NewService(authRepo, jwtSecret, jwtTTL, bcryptCost)
	listService := todo.NewTodoListService(listRepo, log.With("component", "todo_list_service"))
	itemService := todo.NewTodoItemService(itemRepo, listRepo, log.With("component", "todo_item_service"))

	services := &handler.Service{
		Auth:     authService,
		TodoList: listService,
		TodoItem: itemService,
	}

	handlers := handler.NewHandler(services)
	router := handlers.InitRoutes(log)

	port := viper.GetString("port")
	if port == "" {
		port = "8000"
	}

	srv := transportServer.NewServer(port, router, log)

	go func() {
		if err := srv.Run(); err != nil {
			log.Error("HTTP server error", sl.Err(err))
		}
	}()

	log.Info("Server started", slog.String("port", port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", sl.Err(err))
	}

	log.Info("Server exited properly")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	viper.AutomaticEnv()

	return nil
}
