package main

import (
	"log"

	"github.com/TTekmii/todo-list-app"
	handler "github.com/TTekmii/todo-list-app/package/handlers"
	"github.com/TTekmii/todo-list-app/package/repository"
	"github.com/TTekmii/todo-list-app/package/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
