package shards

import (
	"context"
	"database/sql"
	"errors"
	"io/fs"

	"github.com/adlio/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mozgio/database"
	"github.com/skamenetskiy/sharding"
)

func Driver(
	dsn []string,
	strategy sharding.Strategy[uint64, *sql.DB],
) database.Driver[sharding.Cluster[uint64, *sql.DB]] {
	return &driver{
		dsn:      dsn,
		strategy: strategy,
	}
}

type driver struct {
	dsn      []string
	strategy sharding.Strategy[uint64, *sql.DB]
	conn     sharding.Cluster[uint64, *sql.DB]
}

func (d *driver) Connect() (sharding.Cluster[uint64, *sql.DB], error) {
	var (
		err     error
		configs = make([]sharding.ShardConfig, len(d.dsn))
	)
	for i, addr := range d.dsn {
		configs[i] = sharding.ShardConfig{
			ID:   int64(i + 1),
			Addr: addr,
		}
	}
	d.conn, err = sharding.Connect[uint64, *sql.DB](sharding.Config[uint64, *sql.DB]{
		Connect: func(ctx context.Context, addr string) (*sql.DB, error) {
			conn, err := sql.Open("mysql", addr)
			if err != nil {
				return nil, errors.Join(err, errFailedToConnect)
			}
			return conn, nil
		},
		Shards:   configs,
		Context:  context.Background(),
		Strategy: d.strategy,
	})
	if err != nil {
		return nil, errors.Join(err, errClusterInit)
	}
	return d.conn, nil
}

func (d *driver) Close() error {
	return d.conn.Each(func(s sharding.Shard[*sql.DB]) error {
		return s.Conn().Close()
	})
}

func (d *driver) Migrate(files fs.FS, pattern string) error {
	migrations, err := schema.FSMigrations(files, pattern)
	if err != nil {
		return errors.Join(err, database.ErrFailedToReadMigrations)
	}
	opts := []schema.Option{
		schema.WithDialect(schema.MySQL),
	}
	return d.conn.Each(func(s sharding.Shard[*sql.DB]) error {
		if err = schema.NewMigrator(opts...).Apply(s.Conn(), migrations); err != nil {
			return errors.Join(err, database.ErrFailedToMigrate)
		}
		return nil
	})
}

var (
	errFailedToConnect = errors.New("failed to connect")
	errClusterInit     = errors.New("failed to init cluster")
)
