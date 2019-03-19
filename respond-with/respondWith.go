package respondWith

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rhallora-heidelberg/handle"
)

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

func StatusCode(status int) handle.Response {
	return handle.Response{
		StatusCode: status,
		Body:       strings.NewReader(http.StatusText(status)),
	}
}

func Stringf(message string, v ...interface{}) handle.Response {
	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}

	return handle.Response{
		StatusCode: http.StatusOK,
		Body:       strings.NewReader(message),
	}
}

func Bytes(b []byte, status int) handle.Response {
	return handle.Response{
		StatusCode: status,
		Body:       bytes.NewReader(b),
	}
}

func JSONOrError(obj interface{}) handle.Response {
	var b []byte
	var err error

	if m, ok := obj.(json.Marshaler); ok {
		b, err = m.MarshalJSON()
	} else {
		b, err = json.Marshal(obj)
	}

	if err != nil {
		return Errorf(http.StatusInternalServerError, "server failed to marshal JSON response")
	}

	return Bytes(b, http.StatusOK).WithHeaderOptions(setContentJSON)
}

func setContentJSON(hdr http.Header) {
	hdr.Set("Content-Type", "application/json")
}
