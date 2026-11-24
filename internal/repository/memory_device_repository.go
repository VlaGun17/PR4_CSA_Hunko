package repository

import (
	"sync"
	"PR4_CSA_Hunko/internal/domain"
	"fmt"
)

type MemoryDeviceRepository struct {
	mu sync.RWMutex
	storage map[string]domain.Device
}

func NewMemoryDeviceRepository() *MemoryDeviceRepository {
	return &MemoryDeviceRepository{
		storage: make(map[string]domain.Device),
	}
}

func (r *MemoryDeviceRepository) GetAll() ([]domain.Device, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	devices := make([]domain.Device, 0, len(r.storage))
	for _, device := range r.storage {
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *MemoryDeviceRepository) GetByID(id string) (domain.Device, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	device, exists := r.storage[id]
	return device, exists
}

func (r *MemoryDeviceRepository) Create(device domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.storage[device.ID]; exists {
		return fmt.Errorf("device with ID %s already exists", device.ID)
	}
	r.storage[device.ID] = device
	return nil
}

func (r *MemoryDeviceRepository) Update(id string, device domain.Device) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.storage[id]; exists {
		r.storage[id] = device
		return true
	}
	return false
}

func(r *MemoryDeviceRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.storage[id]; exists {
		delete(r.storage, id)
		return true
	}
	return false
}