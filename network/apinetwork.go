package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/grimthereaper/lockout/pb"
)

func (api *API) serve() error {
	return http.ListenAndServe(
		fmt.Sprintf("%v:%v", api.host, api.port),
		api.server,
	)
}

func (api *API) registerHandlers() {
	api.server.HandleFunc(`/api/v0/ip/whitelist`, checkIPAPI)
}

func checkIPAPI(rw http.ResponseWriter, r *http.Request, p *regexp.Regexp) {
	// Might as well reuse types.
	var request pb.IPCheckRequest

	// Using decoder is the best idea. Straight unmarshal is bad.
	decoder := json.NewDecoder(r.Body)
	// NOTE/TODO: Can't remember if this is needed.
	defer r.Body.Close()

	decoder.Decode(&request)

	whitelisted, err := checkIP(request.GetIp(), request.GetCountries())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		// NOTE: Don't use anonymous types, make solid type/function.
		serveFormatted(rw, struct {
			Error string `json:"error"`
		}{Error: err.Error()})
		return
	}

	serveFormatted(rw, pb.IPCheckResponse{Whitelisted: whitelisted})
}
