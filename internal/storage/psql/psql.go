package psql

import (
	"context"
	"fmt"

	"javacode-test/internal/models"
	"javacode-test/internal/storage"

	uuid "github.com/jackc/pgx/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	conn *pgxpool.Pool
}

type Wallet struct {
	WalletID uuid.UUID
	Amount   int
}

var (
	dbInitComm = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s UUID PRIMARY KEY, %s INTEGER);",
		models.TableName, models.ColWalletID, models.ColAmount,
	)

	reqSelectAmountById = fmt.Sprintf("SELECT %s FROM %s WHERE %s=$1",
		models.ColAmount, models.TableName, models.ColWalletID,
	)

	reqApplyOperation = fmt.Sprintf("UPDATE %s SET %s=$1 WHERE %s=$2",
		models.TableName, models.ColAmount, models.ColWalletID,
	)
)

func New(dbURL string) (*Storage, error) {
	const op = "storage.psql.New"

	conn, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w, %s", op, err, dbURL)
	}

	_, err = conn.Exec(context.Background(), dbInitComm)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) UpdateAmountByID(ctx context.Context, log *zap.Logger, walletId string, deltaAmount int) error {
	const op = "storage.psql.ApplyOperation"

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	defer tx.Rollback(ctx)

	var origAmount int
	err = tx.QueryRow(ctx, reqSelectAmountById, walletId).Scan(&origAmount)
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	newAmount := origAmount + deltaAmount
	if newAmount < 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrBalance)
	}

	_, err = tx.Exec(ctx, reqApplyOperation, newAmount, walletId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAmountByID(ctx context.Context, log *zap.Logger, walletId string) (int, error) {
	const op = "storage.psql.GetAmount"

	var amount int
	err := s.conn.QueryRow(ctx, reqSelectAmountById, walletId).Scan(&amount)
	if err != nil {
		log.Error(err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return amount, nil
}
