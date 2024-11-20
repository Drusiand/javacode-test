package apply

import (
	"context"
	"errors"
	"net/http"

	"javacode-test/internal/config/postconfig"
	"javacode-test/internal/http-server/handlers"
	"javacode-test/internal/models"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Request struct {
	WalletId  string `json:"walletId" validate:"required,uuid"`
	Operation string `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount    int    `json:"amount" validate:"required,gt=0"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type WalletApplier interface {
	ApplyOperation(context.Context, *zap.Logger, models.ApplyRequest) error
}

var (
	ErrApply = errors.New("failed to apply operation")
)

type ApplierHandler struct {
	handlerFunc func(w http.ResponseWriter, r *http.Request)
}

func New(cfg *postconfig.PostConfig, log *zap.Logger, wa WalletApplier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.applier.New"

		ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
		defer cancel()

		log := log.With(
			zap.String("op", op),
		)

		var req Request
		err := render.DefaultDecoder(r, &req)
		if err != nil {
			log.Error(handlers.StatusERR)
			render.JSON(w, r, Response{Status: handlers.StatusERR, Error: handlers.ErrDecode.Error()})
			return
		}
		log.Info("request body successfully decoded", zap.Any("request", req))

		validate := validator.New(validator.WithRequiredStructEnabled())
		err = validate.Struct(req)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Response{Status: handlers.StatusERR, Error: handlers.ErrValidate.Error()})
			return
		}

		applyReq := models.ApplyRequest{WalletID: req.WalletId, OperationType: req.Operation, Amount: req.Amount}
		err = wa.ApplyOperation(ctx, log, applyReq)
		if err != nil {
			log.Error(err.Error())
			render.JSON(w, r, Response{Status: handlers.StatusERR, Error: ErrApply.Error()})
			return
		}
		log.Info("Operation successfully applied")
		render.JSON(w, r, Response{Status: handlers.StatusOK, Error: handlers.ErrNoError.Error()})
	}
}
