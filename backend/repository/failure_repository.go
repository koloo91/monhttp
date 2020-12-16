package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
)

func InsertFailure(ctx context.Context, tx *sql.Tx, failure model.Failure) error {
	if _, err := tx.ExecContext(ctx, `INSERT INTO failure (id, service_id, reason, created_at) 
											VALUES ($1, $2, $3, $4)`,
		failure.Id, failure.ServiceId, failure.Reason, failure.CreatedAt); err != nil {
		return err
	}
	return nil
}
