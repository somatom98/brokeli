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
	Expenses       []expenseCall
	Incomes        []incomeCall
	Transfers      []transferCall
	Reimbursements []reimbursementCall
}

type expenseCall struct {
	AccountID   uuid.UUID
	Currency    values.Currency
	Amount      decimal.Decimal
	Category    string
	Description string
}

type incomeCall struct {
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

func (m *DispatcherMock) RegisterExpense(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error {
	m.Expenses = append(m.Expenses, expenseCall{
		AccountID:   accountID,
		Currency:    currency,
		Amount:      amount,
		Category:    category,
		Description: description,
	})
	return nil
}

func (m *DispatcherMock) RegisterIncome(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, category, description string) error {
	m.Incomes = append(m.Incomes, incomeCall{
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
6/26/2025 0:00:00,Lunar,,148.00,DKK,,,Groceries,Expense,A&M coffee+bread,,,
6/3/2025 0:00:00,,Lunar,,,123.82,DKK,Investments,Income,Interests,,,
5/29/2025 0:00:00,Intesa San Paolo,Directa,"1,350.00",EUR,1350,EUR,Investments,Transfer,,,,
5/19/2025 0:00:00,,Lunar,,,80,DKK,Home,Expense,Reimbursement dirt,,,`

		tmpFile, err := os.CreateTemp("", "transactions*.csv")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.WriteString(csvContent)
		require.NoError(t, err)
		tmpFile.Close()

		lunarID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Lunar"))
		intesaID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Intesa San Paolo"))
		directaID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Directa"))

		// act
		err = feature.ImportTransactions(context.Background(), tmpFile.Name())

		// assert
		require.NoError(t, err)

		t.Run("expense correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Expenses, 1)
			expense := dispatcher.Expenses[0]
			assert.Equal(t, lunarID, expense.AccountID)
			assert.Equal(t, values.Currency("DKK"), expense.Currency)
			assert.True(t, decimal.NewFromInt(148).Equal(expense.Amount))
			assert.Equal(t, "Groceries", expense.Category)
		})

		t.Run("income correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Incomes, 1)
			income := dispatcher.Incomes[0]
			assert.Equal(t, lunarID, income.AccountID)
			assert.Equal(t, values.Currency("DKK"), income.Currency)
			assert.True(t, decimal.NewFromFloat(123.82).Equal(income.Amount))
			assert.Equal(t, "Investments", income.Category)
		})

		t.Run("transfer correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Transfers, 1)
			transfer := dispatcher.Transfers[0]
			assert.Equal(t, intesaID, transfer.FromAccountID)
			assert.Equal(t, directaID, transfer.ToAccountID)
			assert.Equal(t, values.Currency("EUR"), transfer.FromCurrency)
			assert.Equal(t, values.Currency("EUR"), transfer.ToCurrency)
			assert.True(t, decimal.NewFromInt(1350).Equal(transfer.FromAmount))
			assert.True(t, decimal.NewFromInt(1350).Equal(transfer.ToAmount))
		})

		t.Run("reimbursement correctly registered", func(t *testing.T) {
			require.Len(t, dispatcher.Reimbursements, 1)
			reimbursement := dispatcher.Reimbursements[0]
			assert.Equal(t, lunarID, reimbursement.AccountID)
			assert.Equal(t, "Reimbursement dirt", reimbursement.From)
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
6/26/2025 0:00:00,Lunar,,0,DKK,,,Groceries,Expense,Zero amount,,,`

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
		assert.Empty(t, dispatcher.Expenses)
		assert.Empty(t, dispatcher.Incomes)
		assert.Empty(t, dispatcher.Transfers)
		assert.Empty(t, dispatcher.Reimbursements)
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
	assert.NotEmpty(t, dispatcher.Expenses)
	t.Logf("Imported %d expenses, %d incomes, %d transfers, %d reimbursements",
		len(dispatcher.Expenses), len(dispatcher.Incomes), len(dispatcher.Transfers), len(dispatcher.Reimbursements))
}
