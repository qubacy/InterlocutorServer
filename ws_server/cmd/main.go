package main

import (
	"ilserver/config"
	"ilserver/server"
	"log"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal("Config initialization has been failed with err:", err)
	}

	// ***

	app := server.NewApp()
	if err := app.Run(); err != nil {
		log.Fatal("App startup failed with err:", err)
	} else {
		log.Println("App startup succeed...")
	}

	// ***

	// TODO: чистое завершение программы (?)
}
