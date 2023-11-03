package model

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
)

type BaseModel struct {
	ID int64 `bun:",pk,autoincrement"`
}

func Initialise(db *bun.DB, ctx context.Context) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		panic(err)
	}

	models := []interface{}{
		&User{},
		&Todo{},
	}

	for _, model := range models {
		_, err = tx.NewCreateTable().Model(model).IfNotExists().Exec(ctx)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
}
