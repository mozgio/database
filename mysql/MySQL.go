package MySQL

import (
	"database/sql"
	"errors"
	"io/fs"

	"github.com/adlio/schema"
	_ "github.com/go-sql-driver/mysql"
	Database "github.com/mozgio/database"
)

func Driver(dsn string) Database.Driver[*sql.DB] {
	return &driver{
		dsn: dsn,
	}
}

type driver struct {
	dsn  string
	conn *sql.DB
}

func (d *driver) Connect() (*sql.DB, error) {
	var err error
	d.conn, err = sql.Open("mysql", d.dsn)
	return d.conn, err
}

func (d *driver) Migrate(files fs.FS, pattern string) error {
	migrations, err := schema.FSMigrations(files, pattern)
	if err != nil {
		return errors.Join(err, Database.ErrFailedToReadMigrations)
	}
	opts := []schema.Option{
		schema.WithDialect(schema.MySQL),
	}
	if err = schema.NewMigrator(opts...).Apply(d.conn, migrations); err != nil {
		return errors.Join(err, Database.ErrFailedToMigrate)
	}
	return nil
}
