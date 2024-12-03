package namhttp

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dev-ijtech/nam-experimental"
	"github.com/dev-ijtech/nam-experimental/namsql"
)

func handleDeviceIndex(logger *log.Logger, deviceStore *namsql.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			devices, err := deviceStore.FindDevices()

			if err != nil {
				logger.Printf("handle device get all: %v", err)
			}

			json.NewEncoder(w).Encode(devices)
		})
}

func handleDeviceView(logger *log.Logger, deviceStore *namsql.DeviceStore) http.Handler {
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
				logger.Printf("handle device get by id: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = json.NewEncoder(w).Encode(device)

			if err != nil {
				logger.Printf("handle device get by id: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		})
}

func handleDeviceMake(logger *log.Logger, deviceStore *namsql.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var device nam.Device

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&device)

			if err != nil {
				logger.Printf("handle device make: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			deviceExists := deviceStore.DeviceExists(device.ID)

			if deviceExists {
				logger.Printf("handle device make: device already exists")
				http.Error(w, "device already exists", http.StatusBadRequest)
				return
			}

			err = deviceStore.MakeDevice(&device)

			if err != nil {
				logger.Printf("handle device make: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func handleDeviceDelete(logger *log.Logger, deviceStore *namsql.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id < 0 {
				http.NotFound(w, r)
				return
			}

			deviceExists := deviceStore.DeviceExists(id)

			if !deviceExists {
				logger.Printf("handle device delete: device does not exist")
				return
			}

			err = deviceStore.DeleteDevice(id)

			if err != nil {
				logger.Printf("handle device delete: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func handleDeviceUpdate(logger *log.Logger, deviceStore *namsql.DeviceStore) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))

			if err != nil || id < 0 {
				http.NotFound(w, r)
				return
			}

			deviceExists := deviceStore.DeviceExists(id)

			if !deviceExists {
				logger.Printf("handle device delete: device does not exist")
				return
			}

			var device nam.Device

			decoder := json.NewDecoder(r.Body)
			err = decoder.Decode(&device)

			if err != nil {
				logger.Printf("handle device make: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = deviceStore.UpdateDevice(id, &device)

			if err != nil {
				logger.Printf("handle device update: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}
