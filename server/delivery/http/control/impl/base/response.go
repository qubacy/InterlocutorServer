package base

import (
	"encoding/json"
	"ilserver/delivery/http/control/dto"
	"ilserver/pkg/utility"
	"io"
	"net/http"
)

func writeHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	//...
}

func writeRawError(w http.ResponseWriter, err error) {
	writeError(w, dto.MakeError(
		utility.UnwrapErrorsToLast(err).Error(), err.Error()))
}

// last write...
// -----------------------------------------------------------------------

func writeError(w http.ResponseWriter, errorObj dto.Error) {
	writeHeaders(w)
	errorJson, err := errorObj.ToJson()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error()) // ignore all.
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, errorJson)
	}
}

func writeJsonOk(w http.ResponseWriter, i any) {
	rawJson, err := json.Marshal(i)
	if err != nil {
		writeRawError(w, err)
		return
	}

	writeHeaders(w)
	w.WriteHeader(http.StatusOK)
	w.Write(rawJson)
}

func writeOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
