// Package handle provides a way to make the execution paths of http handlers easier to trace while eliminating the
// possibility of certain mistakes like responding twice to an http request. It does this by providing the Response
// type so that http responses can be treated primarily as return values instead of side-effects, as well as the With
// function for ease of use.
package handle

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// A PostResponseHook runs at the conclusion of a Response. n and err provide (optional) access to the output of
// copying the Response.Body contents to the http.ResponseWriter.
type PostResponseHook func(n int64, err error)

// A HeaderOption acts on the http.Header of the http.ResponseWriter.
type HeaderOption func(hdr http.Header)

// Response describes how to respond to an http request, and any actions (hooks) which should be performed afterward.
type Response struct {
	StatusCode    int
	Body          io.Reader
	Hooks         []PostResponseHook
	HeaderOptions []HeaderOption
}

// WithHooks attaches the given PostResponseHooks to a Response and returns the result.
func (r Response) WithHooks(hooks ...PostResponseHook) Response {
	if r.Hooks == nil {
		r.Hooks = hooks
		return r
	}

	r.Hooks = append(r.Hooks, hooks...)
	return r
}

// WithHeaderOptions attaches the given HeaderOptions to a Response and returns the result.
func (r Response) WithHeaderOptions(opts ...HeaderOption) Response {
	if r.HeaderOptions == nil {
		r.HeaderOptions = opts
		return r
	}

	r.HeaderOptions = append(r.HeaderOptions, opts...)
	return r
}

// With transforms a function with the signature "func(r *http.Request, ps httprouter.Params) Response" into an httprouter.Handle.
func With(f func(r *http.Request, ps httprouter.Params) Response) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var n int64
		var err error

		res := f(r, ps)

		// protect from invalid WriteHeader panics for zero-value Responses. All others are handled as-is
		if res.StatusCode != 0 {
			w.WriteHeader(res.StatusCode)
		}

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
