package handle

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/julienschmidt/httprouter"
)

func readString(t *testing.T, r io.Reader) string {
	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)

	if b == nil {
		return ""
	}

	return string(b)
}

func doRequest(f HandlerFunc) (*http.Response, error) {
	rr := httptest.NewRecorder()
	handler := httprouter.New()
	handler.GET("/", With(f))

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	handler.ServeHTTP(rr, req)

	return rr.Result(), nil
}

func TestWith(t *testing.T) {
	const helloWorld = "Hello, World!"

	handler := func(_ *http.Request, _ httprouter.Params) Response {
		return Response{
			StatusCode: http.StatusOK,
			Body:       strings.NewReader(helloWorld),
		}
	}

	// basic functionality check, ie can we serve and get back a given response on the other side
	res, err := doRequest(handler)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.EqualValues(t, helloWorld, readString(t, res.Body))
	assert.NoError(t, res.Body.Close())

	// empty responses should default to status 200 with no body
	res, err = doRequest(func(r *http.Request, ps httprouter.Params) Response {
		return Response{}
	})
	assert.NoError(t, err)
	assert.EqualValues(t, "", readString(t, res.Body))
	assert.NoError(t, res.Body.Close())
}

func TestResponse_WithHooks(t *testing.T) {
	n := 0
	countHook := func(_ int64, _ error) {
		n += 1
	}

	// handler should call countHook 3 times when the Response is sent, incrementing n by 3
	handler := func(r *http.Request, ps httprouter.Params) Response {
		return Response{}.WithHooks(countHook, countHook, countHook)
	}

	res, err := doRequest(handler)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.NoError(t, res.Body.Close())

	res, err = doRequest(handler)
	assert.NoError(t, err)
	assert.Equal(t, 6, n)
	assert.NoError(t, res.Body.Close())
}

func TestResponse_WithHeaderOptions(t *testing.T) {
	const (
		cType   = "Content-Type"
		appJSON = "application/json"
	)

	setJSON := func(hdr http.Header) {
		hdr.Set(cType, appJSON)
	}

	handler := func(r *http.Request, ps httprouter.Params) Response {
		return Response{}.WithHeaderOptions(setJSON)
	}

	// ensure we can set and retrieve header values
	res, err := doRequest(handler)
	assert.NoError(t, err)
	assert.Equal(t, appJSON, res.Header.Get(cType))
	assert.NoError(t, res.Body.Close())
}
