package account

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

var (
	ErrAccountNotOpened     = errors.New("account_not_opened")
	ErrNegativeOrNullAmount = errors.New("negative_or_null_amount")
)

func (a *Account) Open(
	name string,
	currency values.Currency,
	happenedAt time.Time,
) (evt event_store.Event, err error) {
	if a.State != State_Unopened {
		return nil, nil
	}

	return &events.Opened{
		AccountID:  a.ID,
		Name:       name,
		Currency:   currency,
		HappenedAt: happenedAt,
	}, nil
}

func (a *Account) UpdateName(
	name string,
	happenedAt time.Time,
) (evt event_store.Event, err error) {
	if a.State != State_Opened {
		return nil, ErrAccountNotOpened
	}

	return &events.NameUpdated{
		AccountID:  a.ID,
		Name:       name,
		HappenedAt: happenedAt,
	}, nil
}

func (a *Account) Deposit(
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	user string,
	happenedAt time.Time,
) (evt event_store.Event, err error) {
	if a.State < State_Opened {
		return nil, ErrAccountNotOpened
	}

	if !amount.IsPositive() {
		return nil, ErrNegativeOrNullAmount
	}

	return &events.MoneyDeposited{
		AccountID:   a.ID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
		User:        user,
		HappenedAt:  happenedAt,
	}, nil
}

func (a *Account) Withdraw(
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
	user string,
	happenedAt time.Time,
) (evt event_store.Event, err error) {
	if a.State < State_Opened {
		return nil, ErrAccountNotOpened
	}

	if !amount.IsPositive() {
		return nil, ErrNegativeOrNullAmount
	}

	return &events.MoneyWithdrawn{
		AccountID:   a.ID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
		User:        user,
		HappenedAt:  happenedAt,
	}, nil
}
