package handle

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

type Respond func(w http.ResponseWriter)

func defaultInternalErr(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	logrus.Errorf("nil handle.Respond called: %s\n", debug.Stack())
}

func With(f func(r *http.Request, ps httprouter.Params) Respond, nilGuard ...Respond) httprouter.Handle {
	var internalErr Respond

	if len(nilGuard) > 0 {
		internalErr = nilGuard[0]
	} else {
		internalErr = defaultInternalErr
	}

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if f == nil {
			internalErr(w)
			return
		}

		respond := f(r, ps)

		if respond == nil {
			internalErr(w)
			return
		}

		respond(w)
	}
}
