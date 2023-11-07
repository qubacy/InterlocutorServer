package base

import (
	"ilserver/delivery/http/control/dto"
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
	"io"
	"net/http"
	"strings"
)

type AdminIdentity struct {
	tokenManager token.Manager
	handler      http.HandlerFunc
}

func NewAdminIdentity(tokenManager token.Manager, handler http.HandlerFunc) *AdminIdentity {
	return &AdminIdentity{
		tokenManager: tokenManager,
		handler:      handler,
	}
}

func (self *AdminIdentity) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := self.parseAuthorizationHeader(&r.Header)
	if err != nil {
		writeRawError(w, err)
		return
	}

	// ***

	self.handler(w, r)
}

// private
// -----------------------------------------------------------------------

func (self *AdminIdentity) parseAuthorizationHeader(header *http.Header) error {
	authValue := header.Get("Authorization")

	authParts := strings.Split(authValue, " ")
	if authParts[0] != "Bearer" {
		return ErrInvalidAuthorizationHeader
	}

	err := self.tokenManager.Validate(authParts[1])
	if err != nil {
		return utility.CreateCustomError(self.parseAuthorizationHeader, err)
	}
	return nil
}

func writeRawError(w http.ResponseWriter, err error) {
	writeError(w, dto.MakeError(
		utility.UnwrapErrorsToLast(err).Error(), err.Error()))
}

func writeError(w http.ResponseWriter, errorObj dto.Error) {
	errorJson, err := dto.ErrorToJson(errorObj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error()) // ignore all.
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, errorJson)
	}
}
