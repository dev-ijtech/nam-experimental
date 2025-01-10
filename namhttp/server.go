package namhttp

import (
	"log"
	"net/http"

	"github.com/dev-ijtech/nam-experimental"
)

func NewServer(logger *log.Logger, deviceStore nam.DeviceStore, southboundService nam.SouthboundService) http.Handler {
	var handler http.Handler
	mux := http.NewServeMux()

	addRoutes(mux, logger, deviceStore, southboundService)

	handler = mux

	handler = loggingMiddleware(logger, handler)

	return handler
}
