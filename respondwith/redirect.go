package respondwith

import (
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/rhallora-heidelberg/handle"
)

// Redirect is almost identical to http.Redirect. From that package:
//
// Redirect replies to the request with a redirect to url,
// which may be a path relative to the request path.
//
// The provided code should be in the 3xx range and is usually
// StatusMovedPermanently, StatusFound or StatusSeeOther.
//
// If the Content-Type header has not been set, Redirect sets it
// to "text/html; charset=utf-8" and writes a small HTML body.
// Setting the Content-Type header to any value, including nil,
// disables that behavior.
func Redirect(r *http.Request, targetURL string, code int) handle.Response {
	targetURL = redirectURL(r, targetURL)

	setLocation := func(hdr http.Header) {
		hdr.Set("Location", hexEscapeNonASCII(targetURL))
	}

	res := handle.Response{
		StatusCode: code,
	}.WithHeaderOptions(setLocation)

	if _, hadCT := r.Header["Content-Type"]; hadCT {
		return res
	}

	if r.Method == "GET" || r.Method == "HEAD" {
		setCT := func(hdr http.Header) {
			hdr.Set("Content-Type", "text/html; charset=utf-8")
		}
		res = res.WithHeaderOptions(setCT)
	}

	if r.Method == "GET" {
		body := "<a href=\"" + htmlEscape(targetURL) + "\">" + http.StatusText(code) + "</a>.\n"
		res.Body = strings.NewReader(body)
	}

	return res
}

// copied from http.Redirect
func redirectURL(r *http.Request, targetURL string) string {
	if u, err := url.Parse(targetURL); err == nil {
		// If targetURL was relative, make its path absolute by
		// combining with request path.
		// The client would probably do this for us,
		// but doing it ourselves is more reliable.
		// See RFC 7231, section 7.1.2
		if u.Scheme == "" && u.Host == "" {
			oldpath := r.URL.Path
			if oldpath == "" { // should not happen, but avoid a crash if it does
				oldpath = "/"
			}

			// no leading http://server
			if targetURL == "" || targetURL[0] != '/' {
				// make relative path absolute
				olddir, _ := path.Split(oldpath)
				targetURL = olddir + targetURL
			}

			var query string
			if i := strings.Index(targetURL, "?"); i != -1 {
				targetURL, query = targetURL[:i], targetURL[i:]
			}

			// clean up but preserve trailing slash
			trailing := strings.HasSuffix(targetURL, "/")
			targetURL = path.Clean(targetURL)
			if trailing && !strings.HasSuffix(targetURL, "/") {
				targetURL += "/"
			}
			targetURL += query
		}
	}

	return targetURL
}

// copied from http.hexEscapeNonASCII
func hexEscapeNonASCII(s string) string {
	newLen := 0
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			newLen += 3
		} else {
			newLen++
		}
	}
	if newLen == len(s) {
		return s
	}
	b := make([]byte, 0, newLen)
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			b = append(b, '%')
			b = strconv.AppendInt(b, int64(s[i]), 16)
		} else {
			b = append(b, s[i])
		}
	}
	return string(b)
}

// copied from http.htmlEscape
func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}

// copied from http.htmlReplacer
var htmlReplacer = strings.NewReplacer( //nolint:gochecknoglobals
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	// "&#34;" is shorter than "&quot;".
	`"`, "&#34;",
	// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	"'", "&#39;",
)
