package control

import (
	"encoding/json"
	"ilserver/service/control"
	"io"
	"net/http"
)

type Handler struct {
	authService  *control.AuthService
	topicService *control.TopicService
}

func NewHandler() *Handler {
	return &Handler{
		authService:  control.NewAuthService(),
		topicService: control.NewTopicService(),
	}
}

// public
// -----------------------------------------------------------------------

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

	err, dtoRes := h.authService.SignIn(dtoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// ***

	resultBytes, err := json.Marshal(dtoRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resultBytes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// -----------------------------------------------------------------------

func (h *Handler) Topic(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.postTopic(w, r)
	}
}

func (h *Handler) Topics(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.postTopics(w, r)
	} else if r.Method == http.MethodGet {
		h.getTopics(w, r)
	}
}

// private
// -----------------------------------------------------------------------

func (h *Handler) postTopic(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	// ***

	jsonTopic, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var dtoReq controlDto.PostTopicReq
	err = json.Unmarshal(jsonTopic, &dtoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ***

	err, dtoRes := h.topicService.PostTopic(dtoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// ***

	resultBytes, err := json.Marshal(dtoRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resultBytes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) postTopics(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) getTopics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	// ***

	jsonBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var dtoReq controlDto.GetTopicsReq
	err = json.Unmarshal(jsonBody, &dtoReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ***

	err, dtoRes := h.topicService.GetTopics(dtoReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// ***

	resultBytes, err := json.Marshal(dtoRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resultBytes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// private
// -----------------------------------------------------------------------
