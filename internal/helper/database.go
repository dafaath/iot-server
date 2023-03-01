package helper

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
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
