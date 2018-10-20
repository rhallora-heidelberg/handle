package respondWith

import (
	"encoding/json"
	"fmt"
	"github.com/rhallora-heidelberg/handle"
	"net/http"
)

func Errorf(status int, message string, v ...interface{}) handle.Respond {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	return func(w http.ResponseWriter) {

		http.Error(w, message, status)
	}
}

func StatusCode(status int) handle.Respond {
	return func(w http.ResponseWriter) {
		w.WriteHeader(status)
		fmt.Fprint(w, http.StatusText(status))
	}
}

func Stringf(message string, v ...interface{}) handle.Respond {
	return func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, message, v...)
	}
}

func Bytes(b []byte, status int) handle.Respond {
	return func(w http.ResponseWriter) {
		w.WriteHeader(status)
		w.Write(b)
	}
}

func JSONOrError(obj interface{}, errorHandler handle.Respond) handle.Respond {
	b, err := json.Marshal(obj)
	if err != nil {
		return errorHandler
	}

	return Bytes(b, http.StatusOK)
}
