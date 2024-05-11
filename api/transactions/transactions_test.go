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

func TestGetAllTransaction(t *testing.T) {
	t.Run("get all transaction successfully", func(t *testing.T) {
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

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

}

func TestCreateTransaction(t *testing.T) {
	t.Run("create transaction successfully", func(t *testing.T) {
		e := echo.New()
		defer e.Close()
		stub := StubTransaction{
			transaction: Transaction{
				Date:      time.Now(),
				Amount:    100,
				Category:  "Food",
				Note:      "Test",
				SpenderID: 1,
			},
		}
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
