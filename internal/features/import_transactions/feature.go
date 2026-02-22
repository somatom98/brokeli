package import_transactions

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

var ErrEmptyAmount = errors.New("empty amounts")

type TransactionDispatcher interface {
	RegisterDeposit(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error
	RegisterWithdrawal(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error
	RegisterTransfer(ctx context.Context, id uuid.UUID, fromAccountID uuid.UUID, fromCurrency values.Currency, fromAmount decimal.Decimal, toAccountID uuid.UUID, toCurrency values.Currency, toAmount decimal.Decimal, category, description string) error
	RegisterReimbursement(ctx context.Context, id uuid.UUID, accountID uuid.UUID, from string, currency values.Currency, amount decimal.Decimal) error
}

type Feature struct {
	httpHandler *http.ServeMux
	dispatcher  TransactionDispatcher
}

func New(
	httpHandler *http.ServeMux,
	api TransactionDispatcher,
) *Feature {
	return &Feature{
		httpHandler: httpHandler,
		dispatcher:  api,
	}
}

func (f *Feature) ImportTransactions(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Some fields might have missing quotes or extra spaces, LazyQuotes handles it.
	reader.LazyQuotes = true

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record: %w", err)
		}

		t, err := newFromRecord(record)
		if errors.Is(err, ErrEmptyAmount) {
			continue
		}
		if err != nil {
			return fmt.Errorf("failed to parse record %v: %w", record, err)
		}

		trxType, err := t.Type()
		if err != nil {
			return fmt.Errorf("failed to get transaction type for record %v: %w", record, err)
		}

		switch trxType {
		case values.TransactionType_Transfer:
			err = f.dispatcher.RegisterTransfer(
				ctx,
				uuid.Must(uuid.NewV7()),
				t.debit.AccountID,
				t.debit.Currency,
				t.debit.Amount,
				t.credit.AccountID,
				t.credit.Currency,
				t.credit.Amount,
				t.category,
				t.description,
			)
			if err != nil {
				return fmt.Errorf("failed to register transfer: %w", err)
			}
		case values.TransactionType_Income:
			err = f.dispatcher.RegisterDeposit(
				ctx,
				uuid.Must(uuid.NewV7()),
				t.credit.AccountID,
				t.credit.Currency,
				t.credit.Amount,
				t.category,
				t.description,
			)
			if err != nil {
				return fmt.Errorf("failed to register income: %w", err)
			}
		case values.TransactionType_Reimbursement:
			fromStr := "unknown"
			if t.description != "" {
				fromStr = t.description
			}

			err = f.dispatcher.RegisterReimbursement(
				ctx,
				uuid.Must(uuid.NewV7()),
				t.credit.AccountID,
				fromStr,
				t.credit.Currency,
				t.credit.Amount,
			)
			if err != nil {
				return fmt.Errorf("failed to register reimbursement: %w", err)
			}
		case values.TransactionType_Expense:
			err = f.dispatcher.RegisterWithdrawal(
				ctx,
				uuid.Must(uuid.NewV7()),
				t.debit.AccountID,
				t.debit.Currency,
				t.debit.Amount,
				t.category,
				t.description,
			)
			if err != nil {
				return fmt.Errorf("failed to register expense: %w", err)
			}
		}
	}

	return nil
}

type transaction struct {
	debit       values.Entry
	credit      values.Entry
	category    string
	trxType     string
	description string
}

func (t transaction) String() string {
	return fmt.Sprintf("debit: %s, credit: %s, type: %s", t.debit, t.credit, t.trxType)
}

func newFromRecord(record []string) (transaction, error) {
	if len(record) < 10 {
		return transaction{}, fmt.Errorf("invalid record length: %v", len(record))
	}

	debitRaw := strings.ReplaceAll(record[3], ",", "")
	if debitRaw == "" {
		debitRaw = "0"
	}
	debitAmount, err := decimal.NewFromString(debitRaw)
	if err != nil {
		return transaction{}, fmt.Errorf("invalid debit: %s, err: %w", record[3], err)
	}

	creditRaw := strings.ReplaceAll(record[5], ",", "")
	if creditRaw == "" {
		creditRaw = "0"
	}
	creditAmount, err := decimal.NewFromString(creditRaw)
	if err != nil {
		return transaction{}, fmt.Errorf("invalid credit: %s, err: %w", record[5], err)
	}

	// Skip empty transactions
	if debitAmount.IsZero() && creditAmount.IsZero() {
		return transaction{}, ErrEmptyAmount
	}

	return transaction{
		debit: values.Entry{
			AccountID: uuid.NewMD5(uuid.NameSpaceOID, []byte(record[1])),
			Currency:  values.Currency(record[4]),
			Amount:    debitAmount,
			Side:      values.Side_Debit,
		},
		credit: values.Entry{
			AccountID: uuid.NewMD5(uuid.NameSpaceOID, []byte(record[2])),
			Currency:  values.Currency(record[6]),
			Amount:    creditAmount,
			Side:      values.Side_Credit,
		},
		category:    record[7],
		trxType:     record[8],
		description: record[9],
	}, nil
}

func (t transaction) Type() (values.TransactionType, error) {
	switch {
	case t.trxType == "Transfer":
		if t.debit.Amount.IsZero() && t.credit.Amount.IsZero() {
			return values.TransactionType_Expense, fmt.Errorf("null amount")
		}
		if t.debit.Amount.IsNegative() {
			return values.TransactionType_Expense, fmt.Errorf("negative debit: %v", t.debit.Amount)
		}
		if t.credit.Amount.IsNegative() {
			return values.TransactionType_Expense, fmt.Errorf("negative credit: %v", t.debit.Amount)
		}
		if t.debit.AccountID == t.credit.AccountID &&
			t.debit.Currency == t.credit.Currency {
			return values.TransactionType_Expense, fmt.Errorf("debit and credit account are the same: %v", t.debit.AccountID)
		}
		return values.TransactionType_Transfer, nil
	case t.trxType == "Income" ||
		(t.trxType == "" && !t.credit.Amount.IsZero() && t.debit.Amount.IsZero()):
		if !t.credit.Amount.IsPositive() {
			return values.TransactionType_Expense, fmt.Errorf("negative credit: %v", t.debit.Amount)
		}
		return values.TransactionType_Income, nil
	case t.trxType == "Expense" && !t.credit.Amount.IsZero() && t.debit.Amount.IsZero():
		return values.TransactionType_Reimbursement, nil
	case t.trxType == "Expense":
		if !t.debit.Amount.IsPositive() {
			return values.TransactionType_Expense, fmt.Errorf("negative debit: %v", t.debit.Amount)
		}
		return values.TransactionType_Expense, nil
	default:
		return values.TransactionType_Expense, fmt.Errorf("unexpected transaction scenario: %s", t)
	}
}
