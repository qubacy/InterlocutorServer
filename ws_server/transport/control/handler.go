package control

import (
	"ilserver/service/control"
	"ilserver/transport/controlDto"
	"net/http"
)

type Handler struct {
	authService *control.AuthService
}

func NewHandler() *Handler {
	return &Handler{
		authService: control.NewAuthService(),
	}
}

// func errRes(error) []byte {

// }

// -----------------------------------------------------------------------

// bytes <-->
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// ***

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !r.PostForm.Has("login") || !r.PostForm.Has("pass") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ***

	dtoReq := controlDto.SignInReq{
		Login: r.PostForm.Get("login"),
		Pass:  r.PostForm.Get("pass"),
	}

	// ***

	//
	err, dtoRes := h.authService.SignIn(dtoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dtoRes.AccessToken = ""

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) PostTopic(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetTopics(w http.ResponseWriter, r *http.Request) {

}
