package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	nanoid "github.com/matoous/go-nanoid"
)

const (
	remoteAddrHeader = "REMOTE_ADDR"
)

type RequestInfo struct {
	UID     string
	URL     string
	Headers *map[string]string
}

func NewRequestInfo(r *http.Request, extractor HeadersExtractor) (*RequestInfo, error) {
	var headers map[string]string

	if extractor == nil {
		headers = make(map[string]string)
	} else {
		headers = extractor.FromRequest(r)
	}

	uid, err := FetchUID(r)

	if err != nil {
		return nil, errors.New("Failed to retrieve connection uid")
	}

	url := r.URL.String()

	if !r.URL.IsAbs() {
		// See https://github.com/golang/go/issues/28940#issuecomment-441749380
		scheme := "http://"
		if r.TLS != nil {
			scheme = "https://"
		}
		url = fmt.Sprintf("%s%s%s", scheme, r.Host, url)
	}

	return &RequestInfo{UID: uid, Headers: &headers, URL: url}, nil
}

// FetchUID safely extracts uid from `X-Request-ID` header or generates a new one
func FetchUID(r *http.Request) (string, error) {
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		return nanoid.Nanoid()
	}

	return requestID, nil
}

func parseCookies(value string, cookieFilter []string) string {
	if len(cookieFilter) == 0 {
		return value
	}

	filter := make(map[string]bool)
	for _, cookie := range cookieFilter {
		filter[cookie] = true
	}

	result := ""
	cookies := strings.Split(value, ";")
	for _, cookie := range cookies {
		parts := strings.Split(cookie, "=")
		if len(parts) != 2 {
			continue
		}

		if filter[parts[0]] {
			result += cookie + ";"
		}
	}

	return result
}
