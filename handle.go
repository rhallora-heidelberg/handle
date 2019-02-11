package handle

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type PostResponseHook func(n int64, err error)

type HeaderOption func(hdr http.Header)

type Response struct {
	StatusCode    int
	Body          io.Reader
	Hooks         []PostResponseHook
	HeaderOptions []HeaderOption
}

func (r Response) WithHooks(hooks ...PostResponseHook) Response {
	if r.Hooks == nil {
		r.Hooks = hooks
		return r
	}

	r.Hooks = append(r.Hooks, hooks...)
	return r
}

func (r Response) WithHeaderOptions(opts ...HeaderOption) Response {
	if r.HeaderOptions == nil {
		r.HeaderOptions = opts
		return r
	}

	r.HeaderOptions = append(r.HeaderOptions, opts...)
	return r
}

func With(f func(r *http.Request, ps httprouter.Params) Response) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var n int64
		var err error

		res := f(r, ps)

		w.WriteHeader(res.StatusCode)
		for _, opt := range res.HeaderOptions {
			opt(w.Header())
		}

		if res.Body != nil {
			n, err = io.Copy(w, res.Body)
		}

		for _, hook := range res.Hooks {
			hook(n, err)
		}
	}
}
