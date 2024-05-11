package transactions

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllTransaction(t *testing.T) {
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

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

}

func TestCreateTransaction(t *testing.T) {
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
		mock.ExpectQuery("cStmt").WithArgs("HongJot", "hong@jot.ok").WillReturnRows(row)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, db)
		err := h.Create(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
	})
}

func TestUpdateTransaction(t *testing.T) {
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
		mock.ExpectQuery("cStmt").WithArgs("HongJot", "hong@jot.ok").WillReturnRows(row)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, db)
		err := h.Update(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
