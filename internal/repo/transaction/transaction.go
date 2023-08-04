package transaction

import (
	"context"
	"github.com/guneyin/bist-tools/pkg/database"
)

func saveAllToDB(ctx context.Context, ts *Transactions) error {
	for _, item := range ts.Items {
		if item.Import {
			err := saveToDB(ctx, &item)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func saveToDB(ctx context.Context, t *Transaction) error {
	db := database.DB.WithContext(ctx)

	err := db.Create(t).Error
	if err != nil {
		return err
	}

	return nil
}
