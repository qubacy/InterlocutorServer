package control

import (
	"log"
	"net/http"
)

type Logging struct {
	handler http.Handler
}

func NewLogging(handler http.Handler) *Logging {
	return &Logging{
		handler: handler,
	}
}

func (self *Logging) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL, "requested")

	self.handler.ServeHTTP(w, req)
}
