package overWs

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// TODO: очистка отключенных пользователей
func BackgroundUpdateRooms(h *Handler) {
	go func() {
		for {
			select {
			case msg := <-h.RoomService.UpdateRoomMsgs:
				c := h.ProfileIdAndConn[msg.ProfileId]
				c.WriteMessage(websocket.TextMessage, msg.BytesResDto)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-time.After(viper.GetDuration("update_rooms.timeout")):
				h.RoomService.BackgroundUpdateRoomsTick()
			}
		}
	}()
}
