package base

import (
	"context"
	"ilserver/delivery/http/control/dto"
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

func (self *Handler) Mux(pathStart string) *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc(pathStart+"sign-in", self.SignIn)

	// ***

	serveMux.Handle(pathStart+"admin", NewAdminIdentity(self.tokenManager, self.Admin))
	serveMux.Handle(pathStart+"topic", NewAdminIdentity(self.tokenManager, self.Topic))

	return serveMux
}

// auth
// -----------------------------------------------------------------------

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.durationProcess)
	defer cancel()

	defer r.Body.Close()
	if r.Method == http.MethodPost {
		http.NotFound(w, r) // only at the highest level!
		return
	}

	h.postSignIn(ctx, w, r)
}

func (h *Handler) postSignIn(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeRawError(w, err)
		return
	}

	reqDto, err := dto.MakePostSignInReqFromJson(jsonBytes)
	if err != nil {
		writeRawError(w, err)
		return
	}

	serviceOut, err := h.services.SignIn(ctx, reqDto.ToServiceInp())
	if err != nil {
		writeRawError(w, err)
		return
	}

	// are there any more errors?

	// ***

	resDto := dto.MakePostSignInRes(serviceOut)
	writeJsonOk(w, resDto)
}

// admin
// -----------------------------------------------------------------------

func (h *Handler) Admin(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.durationProcess)
	defer cancel()

	defer r.Body.Close()
	if r.Method == http.MethodGet {
		h.getAdmin(ctx, w, r)
	} else if r.Method == http.MethodPost {
		h.postAdmin(ctx, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (h *Handler) getAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	io.Copy(io.Discard, r.Body)

	serviceOut, err := h.services.GetAdmins(ctx)
	if err != nil {
		writeRawError(w, err)
		return
	}

	writeJsonOk(w, dto.MakeGetAdminRes(serviceOut))
}

func (h *Handler) postAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeRawError(w, err)
		return
	}

	// <--- delivery
	reqDto, err := dto.MakePostAdminReqFromJson(jsonBytes)
	if err != nil {
		writeRawError(w, err)
		return
	}

	// services <---
	serviceOut, err := h.services.PostAdmin(ctx, reqDto.ToServiceInp())
	if err != nil {
		writeRawError(w, err)
		return
	}

	// ---> delivery
	writeJsonOk(w, dto.MakePostAdminRes(serviceOut))
}

// topic
// -----------------------------------------------------------------------

func (h *Handler) Topic(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.durationProcess)
	defer cancel()

	defer r.Body.Close()
	if r.Method == http.MethodGet {
		h.getTopic(ctx, w, r)
	} else if r.Method == http.MethodPost {
		h.postTopic(ctx, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (h *Handler) getTopic(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	io.Copy(io.Discard, r.Body)

	serviceOut, err := h.services.GetTopics(ctx)
	if err != nil {
		writeRawError(w, err)
		return
	}

	writeJsonOk(w, dto.MakeGetTopicRes(serviceOut.Topics))
}

func (h *Handler) postTopic(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	jsonBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeRawError(w, err)
		return
	}

	reqDto, err := dto.MakePostTopicReqFromJson(jsonBytes)
	if err != nil {
		writeRawError(w, err)
		return
	}

	_, err = h.services.PostTopics(ctx, reqDto.ToServiceInp())
	if err != nil {
		writeRawError(w, err)
		return
	}

	writeOk(w)
}
