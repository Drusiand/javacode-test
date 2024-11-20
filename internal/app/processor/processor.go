package processor

import (
	"context"
	"fmt"

	"javacode-test/internal/app"
	"javacode-test/internal/models"

	"go.uber.org/zap"
)

type Processor struct {
	storage Storage
}

type Storage interface {
	UpdateAmountByID(ctx context.Context, log *zap.Logger, walletId string, deltaAmount int) error
	GetAmountByID(ctx context.Context, log *zap.Logger, walletId string) (int, error)
}

func New(s Storage) *Processor {
	return &Processor{storage: s}
}

func (p *Processor) ApplyOperation(ctx context.Context, log *zap.Logger, r models.ApplyRequest) error {
	const op = "app.processor.ApplyOperation"

	var delta int
	switch r.OperationType {
	case models.OpDeposit:
		delta = r.Amount
	case models.OpWithdraw:
		delta = -r.Amount
	}

	err := p.storage.UpdateAmountByID(ctx, log, r.WalletID, delta)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s: %w", op, app.ErrApply)
	}

	return nil
}

func (p *Processor) GetAmount(ctx context.Context, log *zap.Logger, walletId string) (int, error) {
	const op = "app.processor.ApplyOperation"

	amount, err := p.storage.GetAmountByID(ctx, log, walletId)
	if err != nil {
		log.Error(err.Error())
		return 0, fmt.Errorf("%s: %w", op, app.ErrGetAmount)
	}

	return amount, nil
}
