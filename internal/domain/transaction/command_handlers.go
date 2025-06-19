package transaction

import (
	"errors"

	"github.com/somatom98/brokeli/internal/domain/transaction/events"
)

var (
	ErrNegativeOrNullAmount    = errors.New("negative_or_null_amount")
	ErrInvalidAccount          = errors.New("invalid_account")
	ErrInvalidAmountOrCurrency = errors.New("invalid_amount_or_currency")
)

func HandleCreateExpense(cmd CreateExpense) (evt events.ExpenseCreated, err error) {
	if !cmd.Amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.ExpenseCreated{
		AccountID:   cmd.AccountID,
		Currency:    cmd.Currency,
		Amount:      cmd.Amount,
		Category:    cmd.Category,
		Description: cmd.Description,
	}, nil
}

func HandleCreateIncome(cmd CreateIncome) (evt events.IncomeCreated, err error) {
	if !cmd.Amount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	return events.IncomeCreated{
		AccountID:   cmd.AccountID,
		Currency:    cmd.Currency,
		Amount:      cmd.Amount,
		Category:    cmd.Category,
		Description: cmd.Description,
	}, nil
}

func HandleCreateTransfer(cmd CreateTransfer) (evt events.TransferCreated, err error) {
	if !cmd.FromAmount.IsPositive() ||
		!cmd.ToAmount.IsPositive() {
		return evt, ErrNegativeOrNullAmount
	}

	if cmd.FromAccountID == cmd.ToAccountID {
		return evt, ErrInvalidAccount
	}

	if cmd.FromCurrency == cmd.ToCurrency {
		if !cmd.FromAmount.Equal(cmd.ToAmount) {
			return evt, ErrInvalidAmountOrCurrency
		}
	}

	return events.TransferCreated{
		FromAccountID: cmd.FromAccountID,
		FromCurrency:  cmd.FromCurrency,
		FromAmount:    cmd.FromAmount,
		ToAccountID:   cmd.ToAccountID,
		ToCurrency:    cmd.ToCurrency,
		ToAmount:      cmd.ToAmount,
		Category:      cmd.Category,
		Description:   cmd.Description,
	}, nil
}
