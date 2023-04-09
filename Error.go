package Database

import (
	"errors"
)

var (
	ErrFailedToReadMigrations = errors.New("failed to read database migrations")
	ErrFailedToMigrate        = errors.New("failed to migrate")
)
