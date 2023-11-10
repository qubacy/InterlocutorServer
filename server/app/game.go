package app

import (
	delivery "ilserver/delivery/ws/game"
	service "ilserver/service/game"
	controlStorage "ilserver/storage/control/impl/sql/sqlite"
	gameStorage "ilserver/storage/game"
	"log"
	"net/http"
)

func runGameServer() (string, *http.ServeMux) {
	controlStorageObj, err := controlStorage.Instance()
	if err != nil {
		log.Fatal("Failed to get instance. Err:", err)
	}

	// dependency injection.

	gameHandler := delivery.NewHandler(
		service.NewService(
			gameStorage.NewStorage(),
			controlStorageObj,
		),
	)

	// result...

	pathStart := "/"
	return pathStart, gameHandler.Mux(pathStart)
}
