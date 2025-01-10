package namhttp

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dev-ijtech/nam-experimental"
)

func handleDeviceIndex(logger *log.Logger, deviceStore nam.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			devices, _, err := deviceStore.FindDevices(r.Context(), nam.DeviceFilter{})

			if err != nil {
				logger.Printf("handle device index: %v", err)
			}

			err = encode(w, http.StatusOK, devices)

			if err != nil {
				logger.Printf("handle device index: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})
}

func handleDeviceView(logger *log.Logger, deviceStore nam.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id < 0 {
				http.NotFound(w, r)
				return
			}

			device, err := deviceStore.FindDeviceByID(r.Context(), id)

			if err != nil {
				logger.Printf("handle device view: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = encode(w, http.StatusOK, device)

			if err != nil {
				logger.Printf("handle device view: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})
}

func handleDeviceCreate(logger *log.Logger, deviceStore nam.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			device, problems, err := decodeValid[nam.Device](r)

			if err != nil {
				logger.Printf("handle device make: %v", err)
				http.Error(w, fmt.Sprintf("%s\n%s", err.Error(), problems.String()), http.StatusBadRequest)
				return
			}

			err = deviceStore.CreateDevice(r.Context(), &device)

			if err != nil {
				logger.Printf("handle device make: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func handleDeviceDelete(logger *log.Logger, deviceStore nam.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id < 0 {
				http.NotFound(w, r)
				return
			}

			err = deviceStore.DeleteDevice(r.Context(), id)

			if err != nil {
				logger.Printf("handle device delete: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func handleDeviceUpdate(logger *log.Logger, deviceStore nam.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id < 0 {
				http.NotFound(w, r)
				return
			}

			update, err := decode[nam.DeviceUpdate](r)

			if err != nil {
				logger.Printf("handle device update: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = deviceStore.UpdateDevice(r.Context(), id, &update)

			if err != nil {
				logger.Printf("handle device update: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}
