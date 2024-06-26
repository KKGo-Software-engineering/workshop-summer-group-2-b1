package spender

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateSpender(t *testing.T) {

	t.Run("create spender succesfully when feature toggle is enable", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "HongJot", "email": "hong@jot.ok"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(cStmt).WithArgs("HongJot", "hong@jot.ok").WillReturnRows(row)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, db)
		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{"id": 1, "name": "HongJot", "email": "hong@jot.ok"}`, rec.Body.String())
	})

	t.Run("create spender failed when feature toggle is disable", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "HongJot", "email": "hong@jot.ok"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		cfg := config.FeatureFlag{EnableCreateSpender: false}

		h := New(cfg, nil)
		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("create spender failed when bad request body", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ bad request body }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, nil)
		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid character")
	})

	t.Run("create spender failed on database (feature toggle is enable) ", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name": "HongJot", "email": "hong@jot.ok"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(cStmt).WithArgs("HongJot", "hong@jot.ok").WillReturnError(assert.AnError)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, db)
		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetAllSpender(t *testing.T) {
	t.Run("get all spender succesfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "HongJot", "hong@jot.ok").
			AddRow(2, "JotHong", "jot@jot.ok")
		mock.ExpectQuery(`SELECT id, name, email FROM spender`).WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[{"id": 1, "name": "HongJot", "email": "hong@jot.ok"},
		{"id": 2, "name": "JotHong", "email": "jot@jot.ok"}]`, rec.Body.String())
	})

	t.Run("get all spender failed on database", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		mock.ExpectQuery(`SELECT id, name, email FROM spender`).WillReturnError(assert.AnError)

		h := New(config.FeatureFlag{}, db)
		err := h.GetAll(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestSpenderTransactionById(t *testing.T) {
	t.Run("get spender succesfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/spenders/1/transactions", nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()
		dt := time.Date(2024, 05, 11, 0, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, dt, 100, "category", "expense", "notes", "url_to_image2", 1).
			AddRow(2, dt, 200, "category", "expense", "notes", "url_to_image2", 1)

		mock.ExpectQuery(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id
		FROM transaction
		WHERE spender_id = $1
		LIMIT $2
		OFFSET $3`).
			WithArgs("1", "10", "1").
			WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.SpenderTransactionById(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"transections": [
			  {
				"id": 1,
				"date": "2024-05-11T00:00:00Z",
				"amount": 100,
				"category": "category",
				"transaction_type": "expense",
				"note": "notes",
				"image_url": "url_to_image2",
				"spender_id": 1
			  },
			  {
				"id": 2,
				"date": "2024-05-11T00:00:00Z",
				"amount": 200,
				"category": "category",
				"transaction_type": "expense",
				"note": "notes",
				"image_url": "url_to_image2",
				"spender_id": 1
			  }
			],
			"summary": {
			  "total_income": 0,
			  "total_expenses": 300,
			  "current_balance": -300
			},
			"pagination": {
			  "current_page": 1,
			  "total_pages": 0,
			  "per_page": 10
			}
		  }`, rec.Body.String())
	})
}

func TestSpenderTransactionByIdSummary(t *testing.T) {

	t.Run("get spender summary succesfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/spenders/1/transactions/summary", nil)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		inRow := sqlmock.NewRows([]string{"sum"}).
			AddRow(300)
		exRow := sqlmock.NewRows([]string{"sum"}).
			AddRow(300)

		mock.ExpectQuery(`SELECT SUM(amount) FROM transaction WHERE spender_id = $1 AND transaction_type = $2`).
			WithArgs(1, "income").
			WillReturnRows(inRow)
		mock.ExpectQuery(`SELECT SUM(amount) FROM transaction WHERE spender_id = $1 AND transaction_type = $2`).
			WithArgs(1, "expense").
			WillReturnRows(exRow)

		h := New(config.FeatureFlag{}, db)
		err := h.SpenderTransactionByIdSummary(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{
			"summary": {
			  "total_income": 300,
			  "total_expenses": 300,
			  "current_balance": 0
			}
		  }`, rec.Body.String())
	})
}
