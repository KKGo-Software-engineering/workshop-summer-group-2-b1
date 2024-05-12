package transactions

import (
	"encoding/json"
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

type StubTransaction struct {
	err          error
	transaction  Transaction
	expectToCall map[string]bool
}

func initStub(trans Transaction, err error) StubTransaction {
	stub := StubTransaction{
		err:          err,
		transaction:  trans,
		expectToCall: map[string]bool{},
	}
	return stub
}

func TestGetAllTransaction(t *testing.T) {
	t.Run("get all transaction successfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		dt := time.Date(2024, 05, 11, 9, 07, 29, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url", "spender_id"}).
			AddRow(1, dt, 100, "category", "expense", "notes", "http://www", 1).
			AddRow(2, dt, 200, "category", "expense", "notes", "http://www", 1)

		mock.ExpectQuery(`SELECT id, date, amount, category, transaction_type, note, image_url, spender_id FROM transaction`).
			WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		err := h.GetAll(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `[
			{
			  "id": 1,
			  "date": "2024-05-11T09:07:29Z",
			  "amount": 100,
			  "category": "category",
			  "transaction_type": "expense",
			  "note": "notes",
			  "image_url": "http://www",
			  "spender_id": 1
			},
		   {
			  "id": 2,
			  "date": "2024-05-11T09:07:29Z",
			  "amount": 200,
			  "category": "category",
			  "transaction_type": "expense",
			  "note": "notes",
			  "image_url": "http://www",
			  "spender_id": 1
			}
		  ]`, rec.Body.String())
	})

}

func TestCreateTransaction(t *testing.T) {
	t.Run("create transaction successfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		dt := time.Date(2024, 05, 11, 9, 07, 29, 0, time.UTC)

		stub := initStub(
			Transaction{
				Date:            dt,
				Amount:          1000,
				Category:        "Food",
				Note:            "Lunch",
				TransactionType: "expense",
				ImageUrl:        "https://example.com/image1.jpg",
				SpenderID:       1,
			},
			nil,
		)
		body, _ := json.Marshal(stub.transaction)

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery(cStmt).
			WithArgs(
				stub.transaction.Date,
				stub.transaction.Amount,
				stub.transaction.Category,
				stub.transaction.TransactionType,
				stub.transaction.Note,
				stub.transaction.ImageUrl,
				stub.transaction.SpenderID,
			).
			WillReturnRows(row)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, db)
		err := h.Create(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.JSONEq(t, `{
			"id": 1,
			"date": "2024-05-11T09:07:29Z",
			"amount": 1000,
			"category": "Food",
			"transaction_type": "expense",
			"note": "Lunch",
			"image_url": "https://example.com/image1.jpg",
			"spender_id": 1
		}`, rec.Body.String())

	})
}

func TestUpdateTransaction(t *testing.T) {
	t.Run("update transaction successfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"name": "HongJot", "email": "hong@jot.ok"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		defer db.Close()

		row := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("cStmt").WithArgs("HongJot", "hong@jot.ok").WillReturnRows(row)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, db)
		err := h.Update(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
