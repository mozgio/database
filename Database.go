package Database

import (
	"io/fs"
)

// Driver interface.
type Driver[Conn any] interface {
	Connect() (Conn, error)
	Close() error
	Migrate(data fs.FS, pattern string) error
}
