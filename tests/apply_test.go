package tests

import (
	"net/http"
	"net/url"
	"testing"

	"javacode-test/internal/http-server/handlers"
	"javacode-test/internal/http-server/handlers/apply"

	"github.com/gavv/httpexpect/v2"
)

type postCase struct {
	walletID      string
	operationType string
	amount        int
}

func TestPostOk(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testCases := []postCase{
		{
			walletID:      "11111111-1111-1111-1111-111111111111",
			operationType: "DEPOSIT",
			amount:        1000,
		},
		{
			walletID:      "11111111-1111-1111-1111-111111111111",
			operationType: "WITHDRAW",
			amount:        1000,
		},
	}

	for _, c := range testCases {
		e.POST("/api/v1/wallet").
			WithJSON(apply.Request{
				WalletId:  c.walletID,
				Operation: c.operationType,
				Amount:    c.amount,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("error").
			ContainsValue(handlers.ErrNoError.Error())
	}
}

func TestPostFail(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testCases := []postCase{
		{
			walletID:      "11111111-1111-1111-1111-111111111112",
			operationType: "DEPOSIT",
			amount:        1000,
		},
		{
			walletID:      "12222222-2222-2222-2222-222222222222",
			operationType: "WITHDRAW",
			amount:        2000,
		},
		{
			walletID:      "33333333-3333-3334-3333-333333333333",
			operationType: "WITHDRAW",
			amount:        3000,
		},
	}

	for _, c := range testCases {
		e.POST("/api/v1/wallet").
			WithJSON(apply.Request{
				WalletId:  c.walletID,
				Operation: c.operationType,
				Amount:    c.amount,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("error").
			ContainsValue(apply.ErrApply.Error())
	}
}

func TestPostBadBody(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testCases := []postCase{
		{
			walletID:      "11111111-1111-1111-1111-11111111111",
			operationType: "DEPOSIT",
			amount:        1000,
		},
		{
			walletID:      "22222222-2222-2222-2222-222222222222",
			operationType: "WITHDRA",
			amount:        2000,
		},
		{
			walletID:      "33333333-3333-3333-3333-333333333333",
			operationType: "WITHDRAW",
			amount:        -2000,
		},
	}

	for _, c := range testCases {
		e.POST("/api/v1/wallet").
			WithJSON(apply.Request{
				WalletId:  c.walletID,
				Operation: c.operationType,
				Amount:    c.amount,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("error").
			ContainsValue(handlers.ErrValidate.Error())
	}
}

func TestPostNotEnoughMoney(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testCases := []postCase{
		{
			walletID:      "11111111-1111-1111-1111-111111111111",
			operationType: "WITHDRAW",
			amount:        1001,
		},
	}

	for _, c := range testCases {
		e.POST("/api/v1/wallet").
			WithJSON(apply.Request{
				WalletId:  c.walletID,
				Operation: c.operationType,
				Amount:    c.amount,
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("error").
			ContainsValue(apply.ErrApply.Error())
	}
}
