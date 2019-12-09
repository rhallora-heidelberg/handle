// Package respondWith defines a few convenience functions for consumers of the handle package. Most are simple
// syntactic sugar or mimic common functions from net/http. TemplateOrError and JSONOrError are provided as a model as
// much as for use, as many use-cases would require more sophisticated error handling.
package respondwith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/rhallora-heidelberg/handle"
)

// Errorf produces a Response which mimics the behavior of http.Error in regards to header settings. It also allows
// formatting directives.
func Errorf(status int, message string, v ...interface{}) handle.Response {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	setHeader := func(header http.Header) {
		header.Set("Content-Type", "text/plain; charset=utf-8")
		header.Set("X-Content-Type-Options", "nosniff")
	}

	return handle.Response{
		StatusCode:    status,
		Body:          strings.NewReader(message),
		HeaderOptions: []handle.HeaderOption{setHeader},
	}
}

// StatusCode produces a Response with the given status code in the header and stringified status in the body.
func StatusCode(status int) handle.Response {
	return handle.Response{
		StatusCode: status,
		Body:       strings.NewReader(http.StatusText(status)),
	}
}

// Stringf produces a Response with status 200 and the given string, with optional formatting directives.
func Stringf(message string, v ...interface{}) handle.Response {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	return handle.Response{
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(message),
	}
}

// Bytes produces a Response with the given status code and reads the byte input as the body. No added behavior, just
// syntactic sugar.
func Bytes(status int, b []byte) handle.Response {
	return handle.Response{
		StatusCode: status,
		Body:       bytes.NewReader(b),
	}
}

// JSONOrError attempts to marshal the input as JSON. If successful, it produces a response with that JSON as the body,
// status 200 OK, and sets the content type to application/json. Otherwise, it produces an error response with status
// 500 and the message "server failed to marshal JSON response".
func JSONOrError(obj interface{}) handle.Response {
	var (
		b   []byte
		err error
	)

	// for json.Marshalers, short-circuit to the appropriate method call. Otherwise call json.Marshal.
	if m, ok := obj.(json.Marshaler); ok {
		b, err = m.MarshalJSON()
	} else {
		b, err = json.Marshal(obj)
	}

	if err != nil {
		return Errorf(http.StatusInternalServerError, "server failed to marshal JSON response")
	}

	return Bytes(http.StatusOK, b).WithHeaderOptions(setContentJSON)
}

func setContentJSON(hdr http.Header) {
	hdr.Set("Content-Type", "application/json")
}

// TemplateOrError attempts to execute the given template. If successful, it produces a response with the result as the
// body and status 200 OK. Otherwise, it produces an error response with status 500 and the message "server failed to
// execute template".
func TemplateOrError(t *template.Template, data interface{}) handle.Response {
	buf := new(bytes.Buffer)

	if t == nil {
		return Errorf(http.StatusInternalServerError, "server failed to execute template (nil template specified)")
	}

	if err := t.Execute(buf, data); err != nil {
		return Errorf(http.StatusInternalServerError, "server failed to execute template")
	}

	return handle.Response{
		StatusCode: http.StatusOK,
		Body:       buf,
	}
}
