package app

import (
	"context"
	delivery "ilserver/delivery/ws/game"
	service "ilserver/service/game"
	controlStorage "ilserver/storage/control/impl/sql/sqlite"
	gameStorage "ilserver/storage/game"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

func runGameServer() (string, *http.ServeMux) {
	controlStorageObj, err := controlStorage.Instance()
	if err != nil {
		log.Fatal("Failed to get instance. Err:", err)
	}

	options := service.Config{
		TimeoutForUpdateRooms: viper.GetDuration("update_rooms.timeout"),
		IntervalFromLastUpdateToNextState: viper.GetDuration(
			"background" +
				".update_room" +
				".with_searching_state" +
				".interval_from_last_update_to_next_state",
		),
		ChattingStageDuration: viper.GetDuration("found_game.chatting_stage_duration"),
		ChoosingStageDuration: viper.GetDuration("found_game.choosing_stage_duration"),
	}

	// dependency injection.

	gameHandler := delivery.NewHandler(
		service.NewService(
			context.TODO(),
			options,
			gameStorage.Instance(),
			controlStorageObj,
		),
	)

	// result...

	pathStart := "/"
	return pathStart, gameHandler.Mux(pathStart)
}
