package app

import (
	"ilserver/config"

	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

func Run() {
	if err := config.Initialize(); err != nil {
		log.Fatalf("Config initialize failed. Err: %v", err.Error())
	}

	// *** create dependencies...

	serveMux := http.NewServeMux()
	serveMux.Handle(runControlServer())
	serveMux.Handle(runGameServer())

	httpServer := &http.Server{
		Addr:           ":" + viper.GetString("port"),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        serveMux,
	}

	// ***

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Server start failed. Err: %v", err.Error())
	}
}
