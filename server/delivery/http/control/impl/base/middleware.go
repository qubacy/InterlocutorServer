package base

import (
	"ilserver/pkg/token"
	"ilserver/pkg/utility"
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

	if len(authParts) != 2 {
		return ErrInvalidAuthorizationHeader
	}
	if authParts[0] != "Bearer" {
		return ErrInvalidAuthorizationHeader
	}

	// some stupid checks?
	if len(authParts[1]) == 0 {
		return ErrInvalidAuthorizationHeader
	}

	err := self.tokenManager.Validate(authParts[1])
	if err != nil {
		return utility.CreateCustomError(self.parseAuthorizationHeader, err)
	}
	return nil
}
