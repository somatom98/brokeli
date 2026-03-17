package events

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

const (
	TypeOpened         string = "AccountOpened"
	TypeNameUpdated    string = "AccountNameUpdated"
	TypeMoneyDeposited string = "MoneyDeposited"
	TypeMoneyWithdrawn string = "MoneyWithdrawn"
)

type Opened struct {
	AccountID  uuid.UUID
	Name       string
	Currency   values.Currency
	HappenedAt time.Time
}

func (e Opened) Type() string {
	return TypeOpened
}

func (e Opened) Content() any {
	return e
}

type NameUpdated struct {
	AccountID  uuid.UUID
	Name       string
	HappenedAt time.Time
}

func (e NameUpdated) Type() string {
	return TypeNameUpdated
}

func (e NameUpdated) Content() any {
	return e
}

type MoneyDeposited struct {
	AccountID  uuid.UUID
	Currency   values.Currency
	Amount     decimal.Decimal
	User       string
	HappenedAt time.Time
}

func (e MoneyDeposited) Type() string {
	return TypeMoneyDeposited
}

func (e MoneyDeposited) Content() any {
	return e
}

type MoneyWithdrawn struct {
	AccountID  uuid.UUID
	Currency   values.Currency
	Amount     decimal.Decimal
	User       string
	HappenedAt time.Time
}

func (e MoneyWithdrawn) Type() string {
	return TypeMoneyWithdrawn
}

func (e MoneyWithdrawn) Content() any {
	return e
}
