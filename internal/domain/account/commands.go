package account

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

var ErrNegativeOrNullAmount = errors.New("negative_or_null_amount")

func (a *Account) Create(createdAt time.Time) (evt events.Created, err error) {
	if a.State >= State_Created {
		return evt, nil
	}

	return events.Created{
		Time: createdAt,
	}, nil
}

func (a *Account) Deposit(user string, currency values.Currency, amount decimal.Decimal, time time.Time) (evt events.MoneyDeposited, err error) {
	if a.State == State_Closed {
		return evt, nil
	}

	if !amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.MoneyDeposited{
		User:     user,
		Currency: currency,
		Amount:   amount,
		Time:     time,
	}, nil
}

func (a *Account) Close(time time.Time) (evt events.AccountClosed, err error) {
	if a.State == State_Closed {
		return evt, nil
	}

	return events.AccountClosed{
		Time: time,
	}, nil
}
