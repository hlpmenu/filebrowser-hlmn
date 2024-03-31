package bolt

import (
	"github.com/asdine/storm/v3"

	"github.com/hlpmenu/filebrowser-hlmn/auth"
	"github.com/hlpmenu/filebrowser-hlmn/settings"
	"github.com/hlpmenu/filebrowser-hlmn/share"
	"github.com/hlpmenu/filebrowser-hlmn/storage"
	"github.com/hlpmenu/filebrowser-hlmn/users"
)

// NewStorage creates a storage.Storage based on Bolt DB.
func NewStorage(db *storm.DB) (*storage.Storage, error) {
	userStore := users.NewStorage(usersBackend{db: db})
	shareStore := share.NewStorage(shareBackend{db: db})
	settingsStore := settings.NewStorage(settingsBackend{db: db})
	authStore := auth.NewStorage(authBackend{db: db}, userStore)

	err := save(db, "version", 2) //nolint:gomnd
	if err != nil {
		return nil, err
	}

	return &storage.Storage{
		Auth:     authStore,
		Users:    userStore,
		Share:    shareStore,
		Settings: settingsStore,
	}, nil
}
