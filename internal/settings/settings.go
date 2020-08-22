package settings

import (
	"context"
	"sync"

	"github.com/mattrax/Mattrax/internal/db"
)

// Service allow safely retrieving and setting of server settings
type Service struct {
	settings     db.Setting
	settingsLock sync.RWMutex
}

// Get safely returns the servers settings
func (s *Service) Get() db.Setting {
	s.settingsLock.RLock()
	var settings = s.settings
	s.settingsLock.RUnlock()
	return settings
}

// New initialises a new settings service
func New(q *db.Queries) (*Service, error) {
	settings, err := q.Settings(context.Background())
	if err != nil {
		return nil, err
	}

	if settings.TenantName == "" {
		settings.TenantName = "Mattrax"
	}

	return &Service{
		settings: settings,
	}, nil
}
