package events

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

const (
	Type_Created        string = "Created"
	Type_MoneyDeposited string = "MoneyDeposited"
	Type_AccountClosed  string = "AccountClosed"
)

type Created struct {
	Time time.Time
}

func (e Created) Type() string {
	return Type_Created
}

func (e Created) Content() any {
	return e
}

type MoneyDeposited struct {
	User     string
	Currency values.Currency
	Amount   decimal.Decimal
	Time     time.Time
}

func (e MoneyDeposited) Type() string {
	return Type_MoneyDeposited
}

func (e MoneyDeposited) Content() any {
	return e
}

type AccountClosed struct {
	Time time.Time
}

func (e AccountClosed) Type() string {
	return Type_AccountClosed
}

func (e AccountClosed) Content() any {
	return e
}
