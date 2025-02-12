package namhttp

import (
	"log"
	"net/http"

	"github.com/dev-ijtech/nam-experimental"
)

func loggingMiddleware(logger *log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s - %s %s %s\n", r.RemoteAddr, r.Proto, r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func addRoutes(mux *http.ServeMux, logger *log.Logger, deviceStore nam.DeviceStore, southboundService nam.SouthboundService) {
	mux.Handle("GET /{$}", http.NotFoundHandler())

	// Register devices handlers
	{
		mux.Handle("GET /devices", handleDeviceIndex(logger, deviceStore))
		mux.Handle("GET /devices/{id}", handleDeviceView(logger, deviceStore))
		mux.Handle("DELETE /devices/{id}", handleDeviceDelete(logger, deviceStore))
		mux.Handle("PATCH /devices/{id}", handleDeviceUpdate(logger, deviceStore))
		mux.Handle("POST /devices", handleDeviceCreate(logger, deviceStore))
		mux.Handle("POST /devices/sync", syncAllDevices(logger, deviceStore, southboundService))
		mux.Handle("POST /devices/{id}/sync", syncDevice(logger, deviceStore, southboundService))
	}
}
