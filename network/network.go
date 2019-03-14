package network

import (
	"google.golang.org/grpc"
)

// Serve will serve this application's network stack on provided host and port.
// WARNING: BLOCK
func Serve(host string, port int, grpc bool) (api *API, err error) {
	api = &API{
		host:   host,
		port:   port,
		server: &RegexpHandler{},
	}

	if grpc {
		return api, api.startGRPC()
	}

	api.registerHandlers()

	return api, api.serve()
}

// API is an abstracted form of our API. Written so you can have more than one!
type API struct {
	host       string
	port       int
	server     *RegexpHandler
	grpcServer *grpc.Server
}
