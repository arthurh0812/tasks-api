package main

import (
	"fmt"
	"log"

	"github.com/arthurh0812/tasks-api/pkg/api/handlers"
	"github.com/arthurh0812/tasks-api/pkg/api/middleware"

	"github.com/kataras/iris/v12"
)

func init() {
	configureEnvVar(envKeys)
	err := configurePwd()
	if err != nil {
		log.Fatalf("init: %v", err)
	}
}

func main() {
	app := iris.Default()

	app.Use(middleware.SetAccessControlHeaders)

	taskHandler := handlers.NewTaskHandler(env)
	taskGroup := app.Party("/")

	taskGroup.Get("/", taskHandler.GetTasks)
	taskGroup.Post("/", taskHandler.CreateTask)

	err := app.Run(iris.Addr(getAddr()))
	if err != nil {
		log.Fatal(err)
	}
}

func getAddr() string {
	port := env["PORT"]
	if len(port) == 0 {
		port = "8080"
	}
	return fmt.Sprintf(":" + port)
}