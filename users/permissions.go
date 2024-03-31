package users

// Permissions describe a user's permissions.
type Permissions struct {
	Admin    bool
	Execute  bool
	Create   bool
	Rename   bool
	Modify   bool
	Delete   bool
	Share    bool
	Download bool
}
