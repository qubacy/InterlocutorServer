package control

import (
	"ilserver/delivery/http/control/impl/base"
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

// pathStart ---> ends with '/'
// serveMux  ---> handler itself.
func (self *Handler) Mux(pathStart string) *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(pathStart+"ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "pong")
	})

	// ***

	baseHandler := base.NewHandler(
		self.services, self.tokenManager,
	)

	pathStart = pathStart + "api/"
	serveMux.Handle(
		pathStart, NewLogging(
			baseHandler.Mux(pathStart)))
	return serveMux
}
