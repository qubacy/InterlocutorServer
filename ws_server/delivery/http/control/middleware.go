package control

import "net/http"

type LoggingMiddleware struct {
	handler http.Handler
}
