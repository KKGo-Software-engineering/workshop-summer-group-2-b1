package transactions

import (
	"database/sql"
	"time"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
)

type Transaction struct {
	ID              int64     `json:"id"`
	Date            time.Time `json:"date"`
	Amount          float64   `json:"amount"`
	Category        string    `json:"category"`
	TransactionType string    `json:"transaction_type"`
	Note            string    `json:"note"`
	ImageUrl        string    `json:"image_url"`
	SpenderID       int64     `json:"spender_id"`
}

type handler struct {
	flag config.FeatureFlag
	db   *sql.DB
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}

func (h handler) GetAll(c echo.Context) error {
	// rows, err := db
	return nil
}

func (h handler) Create(c echo.Context) error {
	return nil
}

func (h handler) Update(c echo.Context) error {
	return nil
}
