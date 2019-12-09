// Package handle aims to simplify http handler execution paths and make them easier to trace. In the process, it makes
// it difficult to commit rare but pernicious mistakes like responding twice to an http request. It does this by
// providing the Response type so that http responses can be treated primarily as return values instead of side-effects,
// as well as the With function for ease of use. In short, this is like a poorman's IO monad for go.
//
// This package is intended to be small and modular. That is, you can decide whether to use this package or the more
// flexible streaming semantics of standard go on a per-route basis. The example directory provides some simple
// usage examples to illustrate usage.
//
// Package respondWith is also provided as an optional way to simplify some common types of responses.
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

// A HandlerFunc takes in an http request and produces a Response.
type HandlerFunc func(*http.Request, httprouter.Params) Response

// With transforms a function with the signature "func(r *http.Request, ps httprouter.Params) Response" into
// an httprouter.Handle.
func With(f HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var (
			n   int64
			err error
		)

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
