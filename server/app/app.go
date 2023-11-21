package app

import (
	"context"
	"ilserver/config"
	controlStorage "ilserver/storage/control/impl/sql/sqlite"
	"os"
	"os/signal"
	"syscall"

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

	ctxForGameServer, stopGameServer :=
		context.WithCancel(context.Background())

	serveMux := http.NewServeMux()
	serveMux.Handle(runControlServer())
	serveMux.Handle(runGameServer(ctxForGameServer))

	httpServer := &http.Server{
		Addr:           ":" + viper.GetString("port"),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        serveMux,
	}

	// *** start

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatalf("Server start failed. Err: %v", err)
		}
	}()
	log.Println("Server started")

	// *** graceful shutdown?

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-quit
	log.Println("Server shutting down")

	ctx, cancel := context.WithTimeout(
		context.Background(), viper.GetDuration("shutdown_timeout"))
	defer cancel()

	// *** closing everything in some order

	stopGameServer()
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server stop failed. Err: %v", err)
	}
	controlStorage.Free()
}
