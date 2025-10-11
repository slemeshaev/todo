package main

import (
	"log"

	"github.com/slemeshaev/todo/internal/handler"
	"github.com/slemeshaev/todo/internal/repository"
	todo "github.com/slemeshaev/todo/internal/server"
	"github.com/slemeshaev/todo/internal/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRouters()); err != nil {
		log.Fatalf("error occured while running http-server: %s", err.Error())
	}
}
