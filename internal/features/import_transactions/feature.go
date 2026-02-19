package import_transactions

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type TransactionDispatcher interface {
	RegisterExpense(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error
	RegisterIncome(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error
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

		if len(record) < 10 {
			continue
		}

		from := record[1]
		to := record[2]
		debitStr := strings.ReplaceAll(record[3], ",", "")
		curD := values.Currency(record[4])
		creditStr := strings.ReplaceAll(record[5], ",", "")
		curC := values.Currency(record[6])
		category := record[7]
		inOut := record[8]
		description := record[9]

		// Skip empty transactions
		if debitStr == "" && creditStr == "" {
			continue
		}

		if inOut == "Transfer" || (debitStr != "" && creditStr != "" && from != "" && to != "") {
			fromAccountID := uuid.NewMD5(uuid.NameSpaceOID, []byte(from))
			toAccountID := uuid.NewMD5(uuid.NameSpaceOID, []byte(to))

			fromAmount, _ := decimal.NewFromString(debitStr)
			toAmount, err := decimal.NewFromString(creditStr)
			if err != nil {
				toAmount = fromAmount
			}

			if toAmount.IsZero() {
				toAmount = fromAmount
			}
			if curC == "" {
				curC = curD
			}
			if curD == "" {
				curD = curC
			}

			if !fromAmount.IsPositive() || !toAmount.IsPositive() {
				continue
			}

			if fromAccountID == toAccountID {
				continue
			}

			err = f.dispatcher.RegisterTransfer(
				ctx,
				uuid.Must(uuid.NewV7()),
				fromAccountID,
				curD,
				fromAmount,
				toAccountID,
				curC,
				toAmount,
				category,
				description,
			)
			if err != nil {
				return fmt.Errorf("failed to register transfer: %w", err)
			}
		} else if inOut == "Income" || (inOut == "" && creditStr != "" && debitStr == "") {
			toAccountID := uuid.NewMD5(uuid.NameSpaceOID, []byte(to))
			amount, _ := decimal.NewFromString(creditStr)

			if !amount.IsPositive() {
				continue
			}

			err = f.dispatcher.RegisterIncome(
				ctx,
				uuid.Must(uuid.NewV7()),
				toAccountID,
				curC,
				amount,
				category,
				description,
			)
			if err != nil {
				return fmt.Errorf("failed to register income: %w", err)
			}
		} else if inOut == "Expense" || (inOut == "" && debitStr != "" && creditStr == "") || inOut == "" {
			if creditStr != "" && debitStr == "" {
				// This is a Reimbursement
				toAccountID := uuid.NewMD5(uuid.NameSpaceOID, []byte(to))
				amount, _ := decimal.NewFromString(creditStr)

				if !amount.IsPositive() {
					continue
				}

				fromStr := "unknown"
				if description != "" {
					fromStr = description
				}

				err = f.dispatcher.RegisterReimbursement(
					ctx,
					uuid.Must(uuid.NewV7()),
					toAccountID,
					fromStr,
					curC,
					amount,
				)
				if err != nil {
					return fmt.Errorf("failed to register reimbursement: %w", err)
				}
			} else if debitStr != "" {
				// Normal expense
				fromAccountID := uuid.NewMD5(uuid.NameSpaceOID, []byte(from))
				amount, _ := decimal.NewFromString(debitStr)

				if !amount.IsPositive() {
					continue
				}

				err = f.dispatcher.RegisterExpense(
					ctx,
					uuid.Must(uuid.NewV7()),
					fromAccountID,
					curD,
					amount,
					category,
					description,
				)
				if err != nil {
					return fmt.Errorf("failed to register expense: %w", err)
				}
			}
		}
	}

	return nil
}
