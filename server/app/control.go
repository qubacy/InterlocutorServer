package app

import (
	delivery "ilserver/delivery/http/control"
	token "ilserver/pkg/token/impl"
	service "ilserver/service/control/impl/base"
	storage "ilserver/storage/control/impl/sql/sqlite"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func runControlServer() (string, *http.ServeMux) {
	tokenManager, err := token.NewManager(
		viper.GetString("control_server.token.secret"))

	if err != nil {
		log.Fatal("Failed to create token manager. Err:", err)
	}

	controlStorage, err := storage.Instance()
	if err != nil {
		log.Fatal("Failed to get instance. Err:", err)
	}

	// dependency injection.

	deps := service.Dependencies{
		AccessTokenTTL: viper.GetDuration("control_server.token.duration"),
		TokenManager:   tokenManager,
		Storage:        controlStorage,
	}

	controlHandler := delivery.NewHandler(
		viper.GetDuration("control_server.exchange.max_duration"),
		tokenManager, service.NewServices(deps),
	)

	// result...

	pathStart := "/control/"
	return pathStart, controlHandler.Mux(pathStart)
}
