package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Perceverance7/recipes/internal/handler"
	"github.com/Perceverance7/recipes/internal/models"
	"github.com/Perceverance7/recipes/internal/repository"
	"github.com/Perceverance7/recipes/internal/service"

	"os"

	"github.com/joho/godotenv"
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

func main(){
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil{
		logrus.Fatalf("config init error: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil{
		logrus.Fatalf("error with loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDb(&repository.Config{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		User: viper.GetString("db.user"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName: viper.GetString("db.dbname"),
		SslMode: viper.GetString("db.sslmode"),
	})
	if err != nil{
		logrus.Fatalf("db connect error %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(models.Server)
	go func(){
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil{
			logrus.Fatalf("error running http server: %s", err.Error())
		}
	}()
	
	logrus.Print("Recipes started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Recipes shutting down")

	if err := srv.Shutdown(context.Background()); err != nil{
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil{
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
