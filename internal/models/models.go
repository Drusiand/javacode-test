package models

type ApplyRequest struct {
	WalletID      string
	OperationType string
	Amount        int
}

type GetRequest struct {
	WalletID string
}

const (
	ColWalletID = "walletId"
	ColAmount   = "amount"
	TableName   = "Wallets"
)
