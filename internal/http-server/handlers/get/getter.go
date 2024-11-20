package get

import (
	"context"
	"errors"
	"net/http"
	"time"

	"javacode-test/internal/config/getconfig"
	"javacode-test/internal/http-server/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Request struct {
	WalletId string `json:"walletId"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Amount int    `json:"amount"`
}

type WalletGetter interface {
	GetAmount(context.Context, *zap.Logger, string) (int, error)
}

var (
	ErrGet = errors.New("failed to get amount")
)

func New(getCfg *getconfig.GetConfig, log *zap.Logger, wg WalletGetter, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.getter.New"

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Millisecond)
		defer cancel()

		log := log.With(
			zap.String("op", op),
		)

		var req Request

		req = Request{WalletId: chi.URLParam(r, param)}
		validate := validator.New(validator.WithRequiredStructEnabled())
		err := validate.Var(req, "required,uuid")
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Response{Status: handlers.StatusERR, Error: handlers.ErrValidate.Error()})
			return
		}

		amount, err := wg.GetAmount(ctx, log, req.WalletId)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Response{Status: handlers.StatusERR, Error: ErrGet.Error()})
			return
		}
		log.Info("amout successfully retrieved")
		render.JSON(w, r, Response{Status: "OK", Error: "No error", Amount: amount})
	}
}
