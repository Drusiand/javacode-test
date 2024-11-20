package tests

import (
	"net/http"
	"net/url"
	"testing"

	"javacode-test/internal/http-server/handlers"
	"javacode-test/internal/http-server/handlers/get"

	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8080"
)

type getCase struct {
	walletID string
	amount   int
}

func TestGetAmountOk(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testCases := []getCase{
		{
			walletID: "11111111-1111-1111-1111-111111111111",
			amount:   1000,
		},
		{
			walletID: "22222222-2222-2222-2222-222222222222",
			amount:   2000,
		},
		{
			walletID: "33333333-3333-3333-3333-333333333333",
			amount:   3000,
		},
		{
			walletID: "44444444-4444-4444-4444-444444444444",
			amount:   4000,
		},
		{
			walletID: "55555555-5555-5555-5555-555555555555",
			amount:   5000,
		},
	}

	for _, c := range testCases {
		e.GET("/api/v1/wallets/{walletId}", c.walletID).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("amount").
			ContainsValue(c.amount).
			ContainsKey("error").
			ContainsValue(handlers.ErrNoError.Error())
	}
}

func TestGetAmountFail(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	testCases := []getCase{
		{
			walletID: "11111111-1111-1111-1111-1111111111il",
			amount:   1001,
		},
		{
			walletID: "22222222-2222-2222-2222-222222222221",
			amount:   2000,
		},
		{
			walletID: "33333333-3333-3333-3333-33333333333",
			amount:   3000,
		},
		{
			walletID: "hello there",
			amount:   4000,
		},
		{
			walletID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			amount:   5000,
		},
	}

	for _, c := range testCases {
		e.GET("/api/v1/wallets/{walletId}", c.walletID).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsKey("amount").
			NotContainsValue(c.amount).
			ContainsKey("error").
			ContainsValue(get.ErrGet.Error())
	}
}
