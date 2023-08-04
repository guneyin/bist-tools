package transaction

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guneyin/bist-tools/internal/repo"
	"github.com/guneyin/bist-tools/pkg/database"
	"github.com/guneyin/bist-tools/pkg/redis"
	broker "github.com/guneyin/gobist-broker/entity"
	"time"
)

type Transaction struct {
	database.Model
	Hash       string    `json:"hash" gorm:"hash;index;not null" redis:"-"`
	UserID     string    `json:"user_id" gorm:"user_id;index;not null" redis:"-"`
	Broker     string    `json:"broker" gorm:"broker" redis:"broker"`
	Symbol     string    `json:"symbol" gorm:"symbol;not null" redis:"symbol"`
	Date       time.Time `json:"date" gorm:"date;not null" redis:"date"`
	Quantity   int       `json:"quantity" gorm:"quantity;not null" redis:"quantity"`
	Price      float64   `json:"price" gorm:"price;not null" redis:"price"`
	TypeCode   int       `json:"type_code" gorm:"type_code;not null" redis:"type_code"`
	Import     bool      `json:"import" gorm:"-" redis:"-"`
	Duplicated bool      `json:"duplicated" gorm:"-" redis:"-"`
}

type Transactions struct {
	key   string
	Items []Transaction `json:"items" redis:"items"`
}

func NewTransactions(key string) *Transactions {
	return &Transactions{
		key:   key,
		Items: make([]Transaction, 0),
	}
}

func NewTransaction(uid uuid.UUID, broker string) *Transaction {
	t := &Transaction{
		UserID: uid.String(),
		Broker: broker,
	}
	t.ID = uuid.New()

	return t
}

var _ repo.AdapterFull = (*Transactions)(nil)

// Transaction

func (t *Transaction) getFromCache(ctx context.Context, key, hash string) {
	d, err := redis.GetH(ctx, key, hash)
	if err != nil {
		return
	}

	t.fromJSON(d)
}

func (t *Transaction) setHash() {
	var checksum int64 = 0

	if t.Duplicated {
		checksum = time.Now().UnixNano()
	}

	s := fmt.Sprintf("%s%s%d%d%f%d%d", t.Broker, t.Symbol, t.Date.Unix(), t.Quantity, t.Price, t.TypeCode, checksum)

	h := md5.New()
	h.Write([]byte(s))

	bs := h.Sum(nil)
	t.Hash = fmt.Sprintf("%x", bs)
}

func (t *Transaction) FromImport(in *broker.Transaction) {
	t.Symbol = in.Symbol
	t.Date = in.Date
	t.Quantity = in.Quantity
	t.Price = in.Price
	t.TypeCode = int(in.Type)

	t.setHash()
}

func (t *Transaction) toJSON() []byte {
	r, _ := json.Marshal(&t)

	return r
}

func (t *Transaction) fromJSON(d []byte) {
	_ = json.Unmarshal(d, t)
}

// Transactions

func (ts *Transactions) Add(t *Transaction) {
	ts.Items = append(ts.Items, *t)
}

func (ts *Transactions) FromJSON(d []byte) error {
	return json.Unmarshal(d, &ts)
}

func (ts *Transactions) ToJSON() []byte {
	d, _ := json.MarshalIndent(&ts, "", " ")

	return d
}

func (ts *Transactions) GetAddress() string {
	return ts.key
}

func (ts *Transactions) SaveToCache(ctx context.Context) error {
	for i, _ := range ts.Items {
		item := &ts.Items[i]

		key := ts.GetAddress()

		t := new(Transaction)
		t.getFromCache(ctx, key, item.Hash)

		item.Duplicated = item.Symbol == t.Symbol
		item.Import = !item.Duplicated
		item.setHash()

		err := redis.SetH(ctx, key, map[string]any{item.Hash: item.toJSON()})
		if err != nil {
			return err
		}
	}

	return nil
}

func (ts *Transactions) SaveToDB(c *fiber.Ctx) error {
	err := saveAllToDB(c.Context(), ts)
	if err != nil {
		return err
	}

	return ts.DeleteFromCache(c.Context())
}

func (ts *Transactions) GetFromCache(ctx context.Context) error {
	key := ts.GetAddress()

	data, err := redis.GetHAll(ctx, key)
	if err != nil {
		return err
	}

	t := new(Transaction)
	for _, item := range data {
		t.fromJSON([]byte(item))

		ts.Add(t)
	}

	return nil
}

func (ts *Transactions) DeleteFromCache(ctx context.Context) error {
	key := ts.GetAddress()

	return redis.Delete(ctx, key)
}
