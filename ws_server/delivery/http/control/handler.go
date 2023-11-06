package control

import (
	"ilserver/pkg/token"
	service "ilserver/service/control"
	"io"
	"net/http"
)

type Handler struct {
	services     service.Services
	tokenManager token.Manager
}

func NewHandler(
	services service.Services,
	tokenManager token.Manager,
) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

// -----------------------------------------------------------------------

func (self *Handler) Initialize(serveMux *http.ServeMux) {
	serveMux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong")
	})

	//...

}
