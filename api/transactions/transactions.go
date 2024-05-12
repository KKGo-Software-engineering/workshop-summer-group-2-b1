package transactions

import (
	"database/sql"
	"net/http"
	"strconv"
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

type Summary struct {
	TotalIncome    float64 `json:"total_income"`
	TotalExpenses  float64 `json:"total_expenses"`
	CurrentBalance float64 `json:"current_balance"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_pages"`
	PerPage     int `json:"per_page"`
}

type T struct {
	Transections []Transaction `json:"transections"`
	Summary      Summary       `json:"summary"`
	Pagination   Pagination    `json:"pagination"`
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
	uStmt = `UPDATE transaction SET date = $1, amount = $2, category = $3, transaction_type = $4, note = $5, image_url = $6, spender_id = $7 WHERE id = $8;`
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
	logger := mlog.L(c)
	ctx := c.Request().Context()
	var t Transaction
	err := c.Bind(&t)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, t.Date, t.Amount, t.Category, t.TransactionType, t.Note, t.ImageUrl, t.SpenderID).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, Transaction{
		ID:              lastInsertId,
		Date:            t.Date,
		Amount:          t.Amount,
		Category:        t.Category,
		TransactionType: t.TransactionType,
		Note:            t.Note,
		ImageUrl:        t.ImageUrl,
		SpenderID:       t.SpenderID,
	})
}

func (h handler) Update(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	id := c.Param("id")
	idi, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var t Transaction
	err = c.Bind(&t)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := h.db.ExecContext(ctx, uStmt, t.Date, t.Amount, t.Category, t.TransactionType, t.Note, t.ImageUrl, t.SpenderID, idi)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if rows == 0 {
		return c.JSON(http.StatusNotFound, "transaction not found")
	}

	t.ID = idi
	return c.JSON(http.StatusOK, t)
}

func (h handler) GetSummary(id int, t_type string) (float64, error) {
	rows := h.db.QueryRow(`SELECT SUM(amount) FROM transaction WHERE spender_id = $1 AND transaction_type = $2`, id, t_type)

	sum := 0.0
	if err := rows.Scan(&sum); err != nil {
		return 0, err
	}

	return sum, nil
}
