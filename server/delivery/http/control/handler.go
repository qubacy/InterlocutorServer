package control

import (
	"ilserver/delivery/http/control/impl/base"
	"ilserver/pkg/token"
	service "ilserver/service/control"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	durationProcess time.Duration
	tokenManager    token.Manager
	services        service.Services
}

func NewHandler(
	durationProcess time.Duration,
	tokenManager token.Manager,
	services service.Services,
) *Handler {
	return &Handler{
		durationProcess: durationProcess,
		tokenManager:    tokenManager,
		services:        services,
	}
}

// -----------------------------------------------------------------------

// pathStart ---> starts and ends with '/'
// serveMux  ---> handler itself.
func (self *Handler) Mux(pathStart string) *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(pathStart+"ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong")
	})

	// ***

	baseHandler := base.NewHandler(
		self.durationProcess,
		self.tokenManager, self.services,
	)

	pathStart = pathStart + "api/"
	serveMux.Handle(
		pathStart, NewLogging(
			baseHandler.Mux(pathStart)))
	return serveMux
}
