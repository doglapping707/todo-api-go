package presenter

import (
	"github.com/doglapping707/todo-api-go/domain"
	"github.com/doglapping707/todo-api-go/usecase"
)

type findAccountBalancePresenter struct{}

func NewFindAccountBalancePresenter() usecase.FindAccountBalancePresenter {
	return findAccountBalancePresenter{}
}

func (a findAccountBalancePresenter) Output(balance domain.Money) usecase.FindAccountBalanceOutput {
	return usecase.FindAccountBalanceOutput{Balance: balance.Float64()}
}
