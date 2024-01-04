package domain

import (
	"context"
	"time"
)

type TransferID string

func (t TransferID) String() string {
	return string(t)
}

type (
	// 送金情報が持つ抽象的な実装
	TransferRepository interface {
		Create(context.Context, Transfer) (Transfer, error)
		FindAll(context.Context) ([]Transfer, error)
		WithTransaction(context.Context, func(context.Context) error) error
	}

	// 送金情報
	Transfer struct {
		id                   TransferID
		accountOriginID      AccountID
		accountDestinationID AccountID
		amount               Money
		createdAt            time.Time
	}
)

// 送金情報を返却する
func NewTransfer(
	ID TransferID,
	accountOriginID AccountID,
	accountDestinationID AccountID,
	amount Money,
	createdAt time.Time,
) Transfer {
	return Transfer{
		id:                   ID,
		accountOriginID:      accountOriginID,
		accountDestinationID: accountDestinationID,
		amount:               amount,
		createdAt:            createdAt,
	}
}

// 送金IDを返却する
func (t Transfer) ID() TransferID {
	return t.id
}

// 送金者アカウントIDを返却する
func (t Transfer) AccountOriginID() AccountID {
	return t.accountOriginID
}

// 送金先アカウントIDを返却する
func (t Transfer) AccountDestinationID() AccountID {
	return t.accountDestinationID
}

// 送金額を返却する
func (t Transfer) Amount() Money {
	return t.amount
}

// 送金情報の作成日時を返却する
func (t Transfer) CreatedAt() time.Time {
	return t.createdAt
}
