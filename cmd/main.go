package main

import (
	"github.com/Eagoker/recipes"
	"github.com/Eagoker/recipes/internal/handler"
	"github.com/Eagoker/recipes/internal/repository"
	"github.com/Eagoker/recipes/internal/service"
	
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

	srv := new(recipes.Server)
	srv.Run(viper.GetString("port"), handlers.InitRoutes())
}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
