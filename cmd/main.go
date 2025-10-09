package main

import (
	"log"

	"github.com/slemeshaev/todo/internal/handler"
	todo "github.com/slemeshaev/todo/internal/service"
)

func main() {
	handlers := new(handler.Handler)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRouters()); err != nil {
		log.Fatalf("error occured while running http-server: %s", err.Error())
	}
}
