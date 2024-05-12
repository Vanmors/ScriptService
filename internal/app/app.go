package app

import (
	"ScriptService/internal/repository"
	"ScriptService/internal/service"
	"ScriptService/internal/transport"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

func Run() error {
	// загружаем файл конфигурации
	viper.SetConfigFile("config/config.yml")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// Получаем значения из конфигурации
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	dbname := viper.GetString("database.dbname")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")

	repos, err := repository.NewRepositories(dbname, username, password, host, port)

	if err != nil {
		return err
	}

	services, err := service.NewServices(repos)
	handler := transport.NewHandler(services)

	http.HandleFunc("/command", handler.CreateCommand)
	http.HandleFunc("/commands", handler.GetAllCommands)

	port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Сервер запущен на порту %s", port)
	return http.ListenAndServe(":"+port, nil)
}
