package spender

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/KKGo-Software-engineering/workshop-summer/api/transactions"
	"github.com/kkgo-software-engineering/workshop/mlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Spender struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SpenderSummary struct {
	Summary transactions.Summary `json:"summary"`
}

type handler struct {
	flag config.FeatureFlag
	db   *sql.DB
}

func New(cfg config.FeatureFlag, db *sql.DB) *handler {
	return &handler{cfg, db}
}

const (
	cStmt = `INSERT INTO spender (name, email) VALUES ($1, $2) RETURNING id;`
)

func (h handler) Create(c echo.Context) error {
	if !h.flag.EnableCreateSpender {
		return c.JSON(http.StatusForbidden, "create new spender feature is disabled")
	}

	logger := mlog.L(c)
	ctx := c.Request().Context()
	var sp Spender
	err := c.Bind(&sp)
	if err != nil {
		logger.Error("bad request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var lastInsertId int64
	err = h.db.QueryRowContext(ctx, cStmt, sp.Name, sp.Email).Scan(&lastInsertId)
	if err != nil {
		logger.Error("query row error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	logger.Info("create successfully", zap.Int64("id", lastInsertId))
	sp.ID = lastInsertId
	return c.JSON(http.StatusCreated, sp)
}

func (h handler) GetAll(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()

	rows, err := h.db.QueryContext(ctx, `SELECT id, name, email FROM spender`)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var sps []Spender
	for rows.Next() {
		var sp Spender
		err := rows.Scan(&sp.ID, &sp.Name, &sp.Email)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		sps = append(sps, sp)
	}

	return c.JSON(http.StatusOK, sps)
}

func (h handler) SpenderTransactionById(c echo.Context) error {
	logger := mlog.L(c)
	ctx := c.Request().Context()
	id := c.Param("id")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	rows, err := h.db.QueryContext(ctx, `SELECT id, date, amount, category, transaction_type, note, image_url, spender_id
	FROM transaction
	WHERE spender_id = $1
	LIMIT $2
	OFFSET $3`, id, limit, page)
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var ts []transactions.Transaction
	for rows.Next() {
		var t transactions.Transaction
		err = rows.Scan(&t.ID, &t.Date, &t.Amount, &t.Category, &t.TransactionType, &t.Note, &t.ImageUrl, &t.SpenderID)
		if err != nil {
			logger.Error("scan error", zap.Error(err))
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		ts = append(ts, t)
	}

	exp := 0.0
	inc := 0.0
	for _, v := range ts {
		if v.TransactionType == "expense" {
			exp += v.Amount
		} else if v.TransactionType == "income" {
			inc += v.Amount
		}
	}

	sum := transactions.Summary{
		TotalExpenses:  exp,
		TotalIncome:    inc,
		CurrentBalance: inc - exp,
	}

	reCP, err := strconv.Atoi(page)
	if err != nil {
		logger.Error("scan error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	reL, err := strconv.Atoi(limit)
	if err != nil {
		logger.Error("scan error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	pg := transactions.Pagination{
		CurrentPage: reCP,
		TotalPage:   0,
		PerPage:     reL,
	}

	ss := transactions.T{
		Transections: ts,
		Summary:      sum,
		Pagination:   pg,
	}

	return c.JSON(http.StatusOK, ss)
}

func (h handler) SpenderTransactionByIdSummary(c echo.Context) error {
	logger := mlog.L(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	hs := transactions.New(h.flag, h.db)
	income, err := hs.GetSummary(id, "income")
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	expense, err := hs.GetSummary(id, "expense")
	if err != nil {
		logger.Error("query error", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	ss := SpenderSummary{
		Summary: transactions.Summary{
			TotalIncome:    income,
			TotalExpenses:  expense,
			CurrentBalance: income - expense,
		},
	}

	return c.JSON(http.StatusOK, ss)
}
