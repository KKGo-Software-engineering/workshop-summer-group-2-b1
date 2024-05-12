package transactions

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
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

const (
	cStmt = `INSERT INTO transaction (date, amount, category, transaction_type, note,image_url, spender_id) VALUES ($1, $2,$3, $4, $5, $6,$7) RETURNING id;`
)

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction`)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction

		if err := rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl, &t.SpenderID); err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		transactions = append(transactions, t)
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h handler) Create(c echo.Context) error {
	return c.JSON(http.StatusCreated, "")
}

func (h handler) Update(c echo.Context) error {
	return c.JSON(http.StatusOK, "updated")
}
