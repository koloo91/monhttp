package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
)

func InsertCheck(ctx context.Context, tx *sql.Tx, check model.Check) error {
	if _, err := tx.ExecContext(ctx, `INSERT INTO "check" (id, service_id, latency_in_ms, ping_time_in_ms, is_failure, created_at) 
											VALUES ($1, $2, $3, $4, $5, $6)`,
		check.Id, check.ServiceId, check.LatencyInMs, check.PingTimeInMs, check.IsFailure, check.CreatedAt); err != nil {
		return err
	}
	return nil
}
