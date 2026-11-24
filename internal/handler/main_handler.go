package handler

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"PR4_CSA_Hunko/internal/domain"
	"PR4_CSA_Hunko/internal/repository"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

var deviceRepository *repository.MemoryDeviceRepository

func NewHandler(deviceRepo *repository.MemoryDeviceRepository) http.Handler {
	deviceRepository = deviceRepo
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	r.Route("/api/v1/devices", func(r chi.Router) {
		r.Get("/", getAllDevicesHandler)
		r.Post("/", createDeviceHandler)
		r.Put("/{id}", putDeviceHandler)
		r.Delete("/{id}", deleteDeviceHandler)
		r.Get("/{id}", getDeviceByIDHandler)
	})
	return r
}

func getAllDevicesHandler(w http.ResponseWriter, r *http.Request) {
	devices, err := deviceRepository.GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve devices", http.StatusInternalServerError)
		return
	}
	devicesJSON, err := json.Marshal(devices)
	if err != nil {
		http.Error(w, "Failed to marshal devices", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(devicesJSON)
}

func createDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var device domain.Device

	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if device.Name == "" || device.IpAddress == "" || device.Type == "" {
		http.Error(w, "Name, IpAddress, and Type are required", http.StatusBadRequest)
		return
	}

	device.ID = uuid.NewString()

	if err := deviceRepository.Create(device); err != nil {
		http.Error(w, "Failed to create device", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(device)
}

func getDeviceByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	device, exists := deviceRepository.GetByID(id)
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}
	responseJSON, err := json.Marshal(device)
	if err != nil {
		http.Error(w, "Failed to marshal device", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

func putDeviceHandler(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, "id")
	var updatedDevice domain.Device
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&updatedDevice); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, exists := deviceRepository.GetByID(id)
	if !exists {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	updatedDevice.ID = id
	updatedDevice.UpdatedAt = time.Now()

	if ok := deviceRepository.Update(id, updatedDevice); !ok {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedDevice)

}

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, "id")

	if ok := deviceRepository.Delete(id); !ok {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}