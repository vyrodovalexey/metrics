package pgstorage

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type PgStorageWithAttributes struct {
	conn *pgx.Conn
}

func (p *PgStorageWithAttributes) NewDatabaseConnection(ctx context.Context, c string) error {
	var err error
	p.conn, err = pgx.Connect(ctx, c)
	return err
}
