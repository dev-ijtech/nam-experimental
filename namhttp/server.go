package namhttp

import (
	"log"
	"net/http"

	"github.com/dev-ijtech/nam-experimental/namsql"
)

func NewServer(logger *log.Logger, deviceStore *namsql.DeviceStore) http.Handler {
	var handler http.Handler
	mux := http.NewServeMux()

	addRoutes(mux, logger, deviceStore)

	handler = mux

	handler = loggingMiddleware(logger, handler)

	return handler
}
