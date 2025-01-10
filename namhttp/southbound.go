package namhttp

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/dev-ijtech/nam-experimental"
)

func syncDevice(logger *log.Logger, deviceStore nam.DeviceStore, southboundService nam.SouthboundService) http.Handler {
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
				logger.Printf("syncDevice: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			handleSyncDevice(r.Context(), logger, deviceStore, southboundService, device)
		})
}

func syncAllDevices(logger *log.Logger, deviceStore nam.DeviceStore, southboundService nam.SouthboundService) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			devices, _, err := deviceStore.FindDevices(r.Context(), nam.DeviceFilter{})

			if err != nil {
				logger.Printf("syncAllDevices: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var wg sync.WaitGroup

			for _, device := range devices {
				wg.Add(1)

				go func(device *nam.Device) {
					defer wg.Done()
					handleSyncDevice(r.Context(), logger, deviceStore, southboundService, device)
				}(device)
			}

			wg.Wait()
		})
}

func handleSyncDevice(
	ctx context.Context,
	logger *log.Logger,
	deviceStore nam.DeviceStore,
	southboundService nam.SouthboundService,
	device *nam.Device) {

	southboundOps, err := southboundService.DeviceFactory(device)
	if err != nil {
		logger.Printf("failed to sync %s: %s\n", device.Name, err.Error())
		return
	}
	updatedDevice, err := southboundOps.GetDeviceDetails(ctx)
	if err != nil {
		logger.Printf("failed to sync %s: %s\n", device.Name, err.Error())
		return
	}

	deviceStore.UpdateDevice(
		ctx,
		updatedDevice.ID,
		&nam.DeviceUpdate{
			Name:           &updatedDevice.Name,
			ManagementIPv4: &updatedDevice.ManagementIPv4,
			Vendor:         &updatedDevice.Vendor,
			Version:        &updatedDevice.Version})
	logger.Printf("synced device: %v \n", *updatedDevice)
}
