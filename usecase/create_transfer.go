package usecase

import (
	"context"
	"time"

	"github.com/doglapping707/todo-api-go/domain"
)

type (
	// CreateTransferUseCase input port
	CreateTransferUseCase interface {
		Execute(context.Context, CreateTransferInput) (CreateTransferOutput, error)
	}

	// CreateTransferInput input data
	CreateTransferInput struct {
		AccountOriginID      string `json:"account_origin_id" validate:"required,uuid4"`
		AccountDestinationID string `json:"account_destination_id" validate:"required,uuid4"`
		Amount               int64  `json:"amount" validate:"gt=0,required"`
	}

	// CreateTransferPresenter output port
	CreateTransferPresenter interface {
		Output(domain.Transfer) CreateTransferOutput
	}

	// CreateTransferOutput output data
	CreateTransferOutput struct {
		ID                   string  `json:"id"`
		AccountOriginID      string  `json:"account_origin_id"`
		AccountDestinationID string  `json:"account_destination_id"`
		Amount               float64 `json:"amount"`
		CreatedAt            string  `json:"created_at"`
	}

	createTransferInteractor struct {
		transferRepo domain.TransferRepository
		accountRepo  domain.AccountRepository
		presenter    CreateTransferPresenter
		ctxTimeout   time.Duration
	}
)

// NewCreateTransferInteractor creates new createTransferInteractor with its dependencies
func NewCreateTransferInteractor(
	transferRepo domain.TransferRepository,
	accountRepo domain.AccountRepository,
	presenter CreateTransferPresenter,
	t time.Duration,
) CreateTransferUseCase {
	return createTransferInteractor{
		transferRepo: transferRepo,
		accountRepo:  accountRepo,
		presenter:    presenter,
		ctxTimeout:   t,
	}
}


func (t createTransferInteractor) Execute(ctx context.Context, input CreateTransferInput) (CreateTransferOutput, error) {
	// タイムアウトを設定したコンテキストを取得する
	ctx, cancel := context.WithTimeout(ctx, t.ctxTimeout)
	defer cancel()

	var (
		transfer domain.Transfer
		err      error
	)

	err = t.transferRepo.WithTransaction(ctx, func(ctxTx context.Context) error {
		// 残高の更新をする
		if err = t.process(ctxTx, input); err != nil {
			return err
		}

		// 送金情報を成形する
		transfer = domain.NewTransfer(
			domain.TransferID(domain.NewUUID()),
			domain.AccountID(input.AccountOriginID),
			domain.AccountID(input.AccountDestinationID),
			domain.Money(input.Amount),
			time.Now(),
		)

		// 送金情報を作成する
		transfer, err = t.transferRepo.Create(ctxTx, transfer)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return t.presenter.Output(domain.Transfer{}), err
	}

	return t.presenter.Output(transfer), nil
}


func (t createTransferInteractor) process(ctx context.Context, input CreateTransferInput) error {
	// 送金者アカウントを取得する
	origin, err := t.accountRepo.FindByID(ctx, domain.AccountID(input.AccountOriginID))
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			return domain.ErrAccountOriginNotFound
		default:
			return err
		}
	}

	// 送金者の残高から送金額分を引き落とす
	if err := origin.Withdraw(domain.Money(input.Amount)); err != nil {
		return err
	}

	// 送金先アカウントを取得する
	destination, err := t.accountRepo.FindByID(ctx, domain.AccountID(input.AccountDestinationID))
	if err != nil {
		switch err {
		case domain.ErrAccountNotFound:
			return domain.ErrAccountDestinationNotFound
		default:
			return err
		}
	}

	// 送金先アカウントの残高に送金額分を入金する
	destination.Deposit(domain.Money(input.Amount))

	// 送金者アカウントの残高を更新する
	if err = t.accountRepo.UpdateBalance(ctx, origin.ID(), origin.Balance()); err != nil {
		return err
	}

	// 送金先アカウントの残高を更新する
	if err = t.accountRepo.UpdateBalance(ctx, destination.ID(), destination.Balance()); err != nil {
		return err
	}

	return nil
}
