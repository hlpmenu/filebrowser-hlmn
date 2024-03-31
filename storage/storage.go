package storage

import (
	"github.com/hlpmenu/filebrowser-hlmn/auth"
	"github.com/hlpmenu/filebrowser-hlmn/settings"
	"github.com/hlpmenu/filebrowser-hlmn/share"
	"github.com/hlpmenu/filebrowser-hlmn/users"
)

// Storage is a storage powered by a Backend which makes the necessary
// verifications when fetching and saving data to ensure consistency.
type Storage struct {
	Users    users.Store
	Share    *share.Storage
	Auth     *auth.Storage
	Settings *settings.Storage
}
