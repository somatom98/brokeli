package import_transactions_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/internal/features/import_transactions"
)

type DispatcherMock struct {
	Withdrawals    []withdrawalCall
	Deposits       []depositCall
	Transfers      []transferCall
	Reimbursements []reimbursementCall
}

type withdrawalCall struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

type depositCall struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

type transferCall struct {
	FromAccountID uuid.UUID
	FromCurrency  values.Currency
	FromAmount    decimal.Decimal
	ToAccountID   uuid.UUID
	ToCurrency    values.Currency
	ToAmount      decimal.Decimal
	Category      string
	Description   string
}

type reimbursementCall struct {
	AccountID uuid.UUID
	From      string
	Currency  values.Currency
	Amount    decimal.Decimal
}

func (m *DispatcherMock) RegisterWithdrawal(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error {
	m.Withdrawals = append(m.Withdrawals, withdrawalCall{
		AccountID:   accountID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
	})
	return nil
}

func (m *DispatcherMock) RegisterDeposit(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error {
	m.Deposits = append(m.Deposits, depositCall{
		AccountID:   accountID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
	})
	return nil
}

func (m *DispatcherMock) RegisterTransfer(ctx context.Context, id uuid.UUID, fromAccountID uuid.UUID, fromCurrency values.Currency, fromAmount decimal.Decimal, toAccountID uuid.UUID, toCurrency values.Currency, toAmount decimal.Decimal, category, description string) error {
	m.Transfers = append(m.Transfers, transferCall{
		FromAccountID: fromAccountID,
		FromCurrency:  fromCurrency,
		FromAmount:    fromAmount,
		ToAccountID:   toAccountID,
		ToCurrency:    toCurrency,
		ToAmount:      toAmount,
		Category:      category,
		Description:   description,
	})
	return nil
}

func (m *DispatcherMock) RegisterReimbursement(ctx context.Context, id uuid.UUID, accountID uuid.UUID, from string, currency values.Currency, amount decimal.Decimal) error {
	m.Reimbursements = append(m.Reimbursements, reimbursementCall{
		AccountID: accountID,
		From:      from,
		Currency:  currency,
		Amount:    amount,
	})
	return nil
}

func TestImportTransactions(t *testing.T) {
	t.Run("successfully import various transaction types", func(t *testing.T) {
		// arrange
		dispatcher := &DispatcherMock{}
		feature := import_transactions.New(nil, dispatcher)

		csvContent := `Date,From,To,Debit,CurD,Credit,CurC,Type,In/Out,Description,Category,Subcategory,EUR
6/26/2025 0:00:00,Account A,,148.00,DKK,,,Groceries,Expense,Coffee and bread,,,
6/3/2025 0:00:00,,Account A,,,123.82,DKK,Investments,Income,Interests,,,
5/29/2025 0:00:00,Bank B,Broker C,"1,350.00",EUR,1350,EUR,Investments,Transfer,,,,
5/19/2025 0:00:00,,Account A,,,80,DKK,Home,Expense,Refund for repairs,,,
5/21/2025 0:00:00,,Account A,,,150,DKK,Bills,,Implicit Income bonus,,,`

		tmpFile, err := os.CreateTemp("", "transactions*.csv")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		require.NoError(t, err)
		tmpFile.Close()

		accountAID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Account A"))
		bankBID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Bank B"))
		brokerCID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Broker C"))

		// act
		err = feature.ImportTransactions(context.Background(), tmpFile.Name())

		// assert
		require.NoError(t, err)

		t.Run("expense correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Withdrawals, 1)
			expense := dispatcher.Withdrawals[0]
			assert.Equal(t, accountAID, expense.AccountID)
			assert.Equal(t, values.Currency("DKK"), expense.Currency)
			assert.True(t, decimal.NewFromInt(148).Equal(expense.Amount))
			assert.Equal(t, "Groceries", expense.Category)
		})

		t.Run("income correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Deposits, 2)

			// 6/3/2025 0:00:00,,Account A,,,123.82,DKK,Investments,Income,Interests,,,
			income1 := dispatcher.Deposits[0]
			assert.Equal(t, accountAID, income1.AccountID)
			assert.Equal(t, values.Currency("DKK"), income1.Currency)
			assert.True(t, decimal.NewFromFloat(123.82).Equal(income1.Amount))
			assert.Equal(t, "Investments", income1.Category)

			// 5/21/2025 0:00:00,,Account A,,,150,DKK,Bills,,Implicit Income bonus,,,
			income2 := dispatcher.Deposits[1]
			assert.Equal(t, accountAID, income2.AccountID)
			assert.Equal(t, values.Currency("DKK"), income2.Currency)
			assert.True(t, decimal.NewFromInt(150).Equal(income2.Amount))
			assert.Equal(t, "Bills", income2.Category)
		})

		t.Run("transfer correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Transfers, 1)
			transfer := dispatcher.Transfers[0]
			assert.Equal(t, bankBID, transfer.FromAccountID)
			assert.Equal(t, brokerCID, transfer.ToAccountID)
			assert.Equal(t, values.Currency("EUR"), transfer.FromCurrency)
			assert.Equal(t, values.Currency("EUR"), transfer.ToCurrency)
			assert.True(t, decimal.NewFromInt(1350).Equal(transfer.FromAmount))
			assert.True(t, decimal.NewFromInt(1350).Equal(transfer.ToAmount))
		})

		t.Run("reimbursement correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Reimbursements, 1)
			reimbursement := dispatcher.Reimbursements[0]
			assert.Equal(t, accountAID, reimbursement.AccountID)
			assert.Equal(t, "Refund for repairs", reimbursement.From)
			assert.Equal(t, values.Currency("DKK"), reimbursement.Currency)
			assert.True(t, decimal.NewFromInt(80).Equal(reimbursement.Amount))
		})
	})

	t.Run("skip empty or invalid transactions", func(t *testing.T) {
		// arrange
		dispatcher := &DispatcherMock{}
		feature := import_transactions.New(nil, dispatcher)

		csvContent := `Date,From,To,Debit,CurD,Credit,CurC,Type,In/Out,Description,Category,Subcategory,EUR
6/26/2025 0:00:00,,,,,,,,,,,,
6/26/2025 0:00:00,Account A,,0,DKK,,,Groceries,Expense,Zero amount,,,`

		tmpFile, err := os.CreateTemp("", "invalid_transactions*.csv")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		require.NoError(t, err)
		tmpFile.Close()

		// act
		err = feature.ImportTransactions(context.Background(), tmpFile.Name())

		// assert
		require.NoError(t, err)
		assert.Empty(t, dispatcher.Withdrawals)
		assert.Empty(t, dispatcher.Deposits)
		assert.Empty(t, dispatcher.Transfers)
		assert.Empty(t, dispatcher.Reimbursements)
	})

	t.Run("fails on unexpected scenarios", func(t *testing.T) {
		tests := []struct {
			name    string
			content string
		}{
			{
				name:    "transfer with same account and currency",
				content: "Date,From,To,Debit,CurD,Credit,CurC,Type,In/Out,Description,Category,Subcategory,EUR\n6/26/2025 0:00:00,Account A,Account A,100,EUR,100,EUR,,Transfer,Same account,,,",
			},
			{
				name:    "negative debit amount",
				content: "Date,From,To,Debit,CurD,Credit,CurC,Type,In/Out,Description,Category,Subcategory,EUR\n6/26/2025 0:00:00,Account A,,-100,EUR,,,Misc,Expense,Negative,,,",
			},
			{
				name:    "unexpected scenario (empty trxType with debit)",
				content: "Date,From,To,Debit,CurD,Credit,CurC,Type,In/Out,Description,Category,Subcategory,EUR\n6/26/2025 0:00:00,Account A,,100,EUR,,,,Unexpected,,,",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				dispatcher := &DispatcherMock{}
				feature := import_transactions.New(nil, dispatcher)

				tmpFile, err := os.CreateTemp("", "invalid_*.csv")
				require.NoError(t, err)
				defer os.Remove(tmpFile.Name())

				_, err = tmpFile.WriteString(tt.content)
				require.NoError(t, err)
				tmpFile.Close()

				err = feature.ImportTransactions(context.Background(), tmpFile.Name())
				assert.Error(t, err)
			})
		}
	})
}

func TestImportTransactions_Integration(t *testing.T) {
	// arrange
	dispatcher := &DispatcherMock{}
	feature := import_transactions.New(nil, dispatcher)
	filePath := "transactions.csv"

	// Skip if file doesn't exist (e.g. in CI environments)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Skip("transactions.csv not found, skipping integration test")
	}

	// act
	err := feature.ImportTransactions(context.Background(), filePath)

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, dispatcher.Withdrawals)
	t.Logf("Imported %d withdrawals, %d deposits, %d transfers, %d reimbursements",
		len(dispatcher.Withdrawals), len(dispatcher.Deposits), len(dispatcher.Transfers), len(dispatcher.Reimbursements))
}
