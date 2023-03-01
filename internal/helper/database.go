package helper

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func CommitOrRollback(ctx context.Context, tx pgx.Tx, err *error) {
	var actualError error
	errPanic := recover()
	if errPanic != nil {
		actualError = errPanic.(error)
	} else {
		actualError = *err
	}

	if actualError != nil {
		errorRollback := tx.Rollback(ctx)
		if errorRollback != nil {
			log.Println("Error in rollback", errorRollback)
		}
	} else {
		errorCommit := tx.Commit(ctx)
		if errorCommit != nil {
			log.Println("Error in commit", errorCommit)
		}
	}

	if errPanic != nil {
		panic(errPanic)
	}
}

// Querier common interface for pgx connection.
// https://github.com/jackc/pgx/issues/1333
type Querier interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}
