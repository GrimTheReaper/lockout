package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

/*
Source: https://stackoverflow.com/a/6565407
Although I modified it a bit.

TODO: Remove this. I thought we needed regexp routes at first, but from the looks of it we won't.
*/

type route struct {
	pattern *regexp.Regexp
	handler func(http.ResponseWriter, *http.Request, *regexp.Regexp)
}

// RegexpHandler was written by "Evan Shaw" over on Stackoverflow
type RegexpHandler struct {
	routes []*route
}

// HandleFunc will append the handler function to the routes
func (h *RegexpHandler) HandleFunc(regex string, handler func(http.ResponseWriter, *http.Request, *regexp.Regexp)) {
	h.routes = append(h.routes, &route{regexp.MustCompile(regex), handler})
}

// ServeHTTP is for http.ListenAndServe
func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler(w, r, route.pattern)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}

func getFirstCaptureGroup(request *http.Request, pattern *regexp.Regexp) string {
	return pattern.FindStringSubmatch(request.URL.Path)[1] // We can assume that it will be there because it was vetted earlier in the code.
}

// NOTE: Since we re-use the PB objects, maybe use https://godoc.org/github.com/golang/protobuf/jsonpb
func serveFormatted(rw http.ResponseWriter, object interface{}) {
	byts, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		rw.Write([]byte(fmt.Sprintf("Failed to encode response: %v", err)))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Write(byts)
}
