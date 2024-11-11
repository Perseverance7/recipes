package main

import (
	"github.com/Perceverance7/recipes/internal/handler"
	"github.com/Perceverance7/recipes/internal/models"
	"github.com/Perceverance7/recipes/internal/repository"
	"github.com/Perceverance7/recipes/internal/service"

	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Recipes App API
// @version 1.0
// description API Server for Recipes App

// @host localhost:8081
// @BasePath: /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("config init error: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error with loading env variables: %s", err.Error())
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Используйте имя сервиса Redis из docker-compose
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("Ошибка подключения к Redis: %v", err)
	}

	logrus.Println("Подключение к Redis установлено")

	db, err := repository.NewPostgresDb(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   viper.GetString("db.dbname"),
		SslMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("db connect error %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos, rdb)
	handlers := handler.NewHandler(services)

	srv := new(models.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("error running HTTP server: %s", err.Error())
		}
	}()

	logrus.Print("Recipes started")

	// Канал для получения сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Recipes shutting down")

	// Контекст с тайм-аутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	// Закрытие подключения к базе данных
	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on DB connection close: %s", err.Error())
	}

	logrus.Print("Shutdown complete")

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
