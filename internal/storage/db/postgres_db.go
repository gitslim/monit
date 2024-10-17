package db

import (
	"context"
	"fmt"

	"github.com/gitslim/monit/internal/logging"
	"github.com/jackc/pgx/v4"
)

func Connect(ctx context.Context, log *logging.Logger, dsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	log.Debug("conn", conn)
	return conn, nil
}
