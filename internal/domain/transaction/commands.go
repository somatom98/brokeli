package transaction

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

var (
	ErrNegativeOrNullAmount    = errors.New("negative_or_null_amount")
	ErrInvalidAccount          = errors.New("invalid_account")
	ErrInvalidAmountOrCurrency = errors.New("invalid_amount_or_currency")
)

func (a *Transaction) SetExpectedReimbursement(
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
) (evt events.ExpectedReimbursementSet, err error) {
	if a.State > State_Created {
		return evt, nil
	}

	if !amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.ExpectedReimbursementSet{
		AccountID: accountID,
		Currency:  currency,
		Amount:    amount,
	}, nil
}

func (a *Transaction) RegisterExpense(
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
) (evt events.MoneySpent, err error) {
	if a.State > State_Created {
		return evt, nil
	}

	if !amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.MoneySpent{
		AccountID:   accountID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
	}, nil
}

func (a *Transaction) RegisterIncome(
	accountID uuid.UUID,
	currency values.Currency,
	amount decimal.Decimal,
	category string,
	description string,
) (evt events.MoneyReceived, err error) {
	if a.State > State_Created {
		return evt, nil
	}

	if !amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.MoneyReceived{
		AccountID:   accountID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
	}, nil
}

func (a *Transaction) RegisterTransfer(
	fromAccountID uuid.UUID,
	fromCurrency values.Currency,
	fromAmount decimal.Decimal,
	toAccountID uuid.UUID,
	toCurrency values.Currency,
	toAmount decimal.Decimal,
	category string,
	description string,
) (evt events.MoneyTransfered, err error) {
	if a.State > State_Created {
		return evt, nil
	}

	if !fromAmount.IsPositive() ||
		!toAmount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	if fromAccountID == toAccountID {
		return evt, ErrInvalidAccount
	}

	if fromCurrency == toCurrency {
		if !fromAmount.Equal(toAmount) {
			return evt, ErrInvalidAmountOrCurrency
		}
	}

	return events.MoneyTransfered{
		FromAccountID: fromAccountID,
		FromCurrency:  fromCurrency,
		FromAmount:    fromAmount,
		ToAccountID:   toAccountID,
		ToCurrency:    toCurrency,
		ToAmount:      toAmount,
		Category:      category,
		Description:   description,
	}, nil
}

func (a *Transaction) RegisterReimbursement(
	accountID uuid.UUID,
	from string,
	currency values.Currency,
	amount decimal.Decimal,
) (evt events.ReimbursementReceived, err error) {
	if a.State > State_Created {
		return evt, nil
	}

	if !amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.ReimbursementReceived{
		AccountID: accountID,
		From:      from,
		Currency:  currency,
		Amount:    amount,
	}, nil
}
