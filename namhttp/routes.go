package namhttp

import (
	"log"
	"net/http"

	"github.com/dev-ijtech/nam-experimental/namsql"
)

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		logger.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL)
	})
}

func addRoutes(mux *http.ServeMux, logger *log.Logger, deviceStore *namsql.DeviceStore) {
	mux.Handle("GET /{$}", http.NotFoundHandler())

	// Register devices handlers
	{
		mux.Handle("GET /devices", handleDeviceIndex(logger, deviceStore))
		mux.Handle("GET /devices/{id}", handleDeviceView(logger, deviceStore))
		mux.Handle("DELETE /devices/{id}", handleDeviceDelete(logger, deviceStore))
		mux.Handle("PUT /devices/{id}", handleDeviceUpdate(logger, deviceStore))
		mux.Handle("POST /devices", handleDeviceMake(logger, deviceStore))
	}
}
