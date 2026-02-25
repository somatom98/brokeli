package account

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

var (
	ErrAccountAlreadyOpened = errors.New("account_already_opened")
	ErrAccountNotOpened     = errors.New("account_not_opened")
	ErrNegativeOrNullAmount = errors.New("negative_or_null_amount")
)

func (a *Account) Open(
	name string,
	currency values.Currency,
) (evt *events.Opened, err error) {
	if a.State != State_Unopened {
		return nil, ErrAccountAlreadyOpened
	}

	return &events.Opened{
		AccountID: a.ID,
		Name:      name,
		Currency:  currency,
	}, nil
}

func (a *Account) UpdateName(
	name string,
) (evt *events.NameUpdated, err error) {
	if a.State != State_Opened {
		return nil, ErrAccountNotOpened
	}

	return &events.NameUpdated{
		Name: name,
	}, nil
}

func (a *Account) Deposit(
	currency values.Currency,
	amount decimal.Decimal,
	user string,
) (evt *events.MoneyDeposited, err error) {
	if a.State < State_Opened {
		return nil, ErrAccountNotOpened
	}

	if !amount.IsPositive() {
		return nil, ErrNegativeOrNullAmount
	}

	return &events.MoneyDeposited{
		AccountID: a.ID,
		Currency:  currency,
		Amount:    amount,
		User:      user,
	}, nil
}

func (a *Account) Withdraw(
	currency values.Currency,
	amount decimal.Decimal,
	user string,
) (evt *events.MoneyWithdrawn, err error) {
	if a.State < State_Opened {
		return nil, ErrAccountNotOpened
	}

	if !amount.IsPositive() {
		return nil, ErrNegativeOrNullAmount
	}

	return &events.MoneyWithdrawn{
		AccountID: a.ID,
		Currency:  currency,
		Amount:    amount,
		User:      user,
	}, nil
}
