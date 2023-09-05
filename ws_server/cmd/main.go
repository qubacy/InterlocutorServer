package main

import (
	"ilserver/config"
	"ilserver/repository"
	"ilserver/server"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal("Config initialization has been failed with err:", err)
	}
	if err := repository.Init(); err != nil {
		log.Fatal("Repository initialization has been failed with err:", err)
	}

	// ***

	app := server.NewApp()

	// ***

	log.Println("Trying to run server...")

	// run only once!
	if err := app.Run(); err != nil {
		log.Fatal("App startup failed with err:", err)
	}
}
