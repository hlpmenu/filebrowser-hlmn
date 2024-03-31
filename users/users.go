package users

import (
	"path/filepath"
	"regexp"

	"github.com/spf13/afero"

	config "github.com/hlpmenu/filebrowser-hlmn/config"
	"github.com/hlpmenu/filebrowser-hlmn/errors"
	"github.com/hlpmenu/filebrowser-hlmn/files"
	"github.com/hlpmenu/filebrowser-hlmn/rules"
)

// ViewMode describes a view mode.
type ViewMode string

const (
	ListViewMode   ViewMode = "list"
	MosaicViewMode ViewMode = "mosaic"
)

var UserPermissions = Permissions{
	Admin:    true,
	Execute:  true,
	Create:   true,
	Rename:   true,
	Modify:   true,
	Delete:   true,
	Share:    true,
	Download: true,
}

// User describes a user.
type User struct {
	ID           uint
	Username     string
	Password     string
	Scope        string
	Locale       string
	LockPassword bool
	ViewMode     ViewMode `json:"viewMode"`
	SingleClick  bool     `json:"singleClick"`
	Perm         Permissions
	Commands     []string      `json:"commands"`
	Sorting      files.Sorting `json:"sorting"`
	Fs           afero.Fs      `json:"-" yaml:"-"`
	Rules        []rules.Rule  `json:"rules"`
	HideDotfiles bool          `json:"hideDotfiles"`
	DateFormat   bool          `json:"dateFormat"`
}

// Define a global User instance using constants from the config package
// Define a global User instance using constants from the config package
var TheUser = &User{
	ID:           config.ID,
	Username:     config.Username,
	Password:     config.Password,
	Scope:        config.Scope,
	Locale:       config.Locale,
	LockPassword: config.LockPassword,
	ViewMode:     ListViewMode, // use one of the ViewMode constants
	SingleClick:  config.SingleClick,
	Perm:         UserPermissions, // replace with the local variable Permissions
	Commands:     []string{},      // initialize Commands as an empty slice
	Sorting:      files.Sorting{}, // replace with an instance of files.Sorting
	Fs:           afero.NewOsFs(), // replace with an instance of afero.Fs
	Rules:        []rules.Rule{},  // replace with an array of instances of rules.Rule
	HideDotfiles: config.HideDotfiles,
	DateFormat:   config.DateFormat,
}

// GetRules implements rules.Provider.
func (u *User) GetRules() []rules.Rule {
	return u.Rules
}

var checkableFields = []string{
	"Username",
	"Password",
	"Scope",
	"ViewMode",
	"Commands",
	"Sorting",
	"Rules",
}

// Clean cleans up a user and verifies if all its fields
// are alright to be saved.
//
//nolint:gocyclo
func (u *User) Clean(baseScope string, fields ...string) error {
	if len(fields) == 0 {
		fields = checkableFields
	}

	for _, field := range fields {
		switch field {
		case "Username":
			if u.Username == "" {
				return errors.ErrEmptyUsername
			}
		case "Password":
			if u.Password == "" {
				return errors.ErrEmptyPassword
			}
		case "ViewMode":
			if u.ViewMode == "" {
				u.ViewMode = ListViewMode
			}
		case "Commands":
			if u.Commands == nil {
				u.Commands = []string{}
			}
		case "Sorting":
			if u.Sorting.By == "" {
				u.Sorting.By = "name"
			}
		case "Rules":
			if u.Rules == nil {
				u.Rules = []rules.Rule{}
			}
		}
	}

	if u.Fs == nil {
		scope := u.Scope
		scope = filepath.Join(baseScope, filepath.Join("/", scope)) //nolint:gocritic
		u.Fs = afero.NewBasePathFs(afero.NewOsFs(), scope)
	}

	return nil
}

// FullPath gets the full path for a user's relative path.
func (u *User) FullPath(path string) string {
	return afero.FullBaseFsPath(u.Fs.(*afero.BasePathFs), path)
}

// CanExecute checks if an user can execute a specific command.
func (u *User) CanExecute(command string) bool {
	if !u.Perm.Execute {
		return false
	}

	for _, cmd := range u.Commands {
		if regexp.MustCompile(cmd).MatchString(command) {
			return true
		}
	}

	return false
}
