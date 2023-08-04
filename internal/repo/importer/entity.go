package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guneyin/bist-tools/internal/middleware"
	"github.com/guneyin/bist-tools/internal/repo"
	"github.com/guneyin/bist-tools/internal/repo/transaction"
	"github.com/guneyin/bist-tools/pkg/redis"
	broker "github.com/guneyin/gobist-broker"
	"github.com/guneyin/gobist-broker/entity"
	"time"
)

var _ repo.AdapterFull = (*ImportSession)(nil)

type ImportSession struct {
	UserID       uuid.UUID
	DateCreated  time.Time
	Transactions *transaction.Transactions
}

func NewImportSession(c *fiber.Ctx) (*ImportSession, error) {
	is := &ImportSession{
		UserID:      middleware.GetUserID(c),
		DateCreated: time.Now(),
	}

	key := is.GetAddress()

	is.Transactions = transaction.NewTransactions(key)

	err := is.GetFromCache(c.Context())
	if err != nil {
		return nil, err
	}

	return is, nil
}

func (is *ImportSession) process(c *fiber.Ctx, b broker.Broker, ts *entity.Transactions) (*transaction.Transactions, error) {
	uid := middleware.GetUserID(c)

	for _, item := range ts.Items {
		t := transaction.NewTransaction(uid, b.Info().Name)
		t.FromImport(&item)

		is.Transactions.Add(t)
	}

	err := is.Transactions.SaveToCache(c.Context())
	if err != nil {
		return nil, err
	}

	return is.Transactions, nil
}

func (is *ImportSession) GetAddress() string {
	return fmt.Sprintf("user:%s:import", is.UserID)
}

func (is *ImportSession) GetFromCache(ctx context.Context) error {
	return is.Transactions.GetFromCache(ctx)
}

func (is *ImportSession) SaveToCache(ctx context.Context) error {
	return is.Transactions.SaveToCache(ctx)
}

func (is *ImportSession) SaveToDB(c *fiber.Ctx) error {
	return is.Transactions.SaveToDB(c)
}

func (is *ImportSession) DeleteFromCache(ctx context.Context) error {
	key := is.GetAddress()

	return redis.Delete(ctx, key)
}

func (is *ImportSession) ToJSON() []byte {
	d, _ := json.MarshalIndent(&is, "", " ")

	return d
}

func (is *ImportSession) FromJSON(d []byte) error {
	return json.Unmarshal(d, &is)
}
