package transaction_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func TestSetExpectedReimbursement(t *testing.T) {
	t.Run("should emit expected reimbursement set event when amount is positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		accountID := uuid.New()
		amount := decimal.NewFromInt(100)

		// act
		evt, err := tx.SetExpectedReimbursement(accountID, values.Currency("USD"), amount)

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.ExpectedReimbursementSet{
			AccountID: accountID,
			Currency:  values.Currency("USD"),
			Amount:    amount,
		}, evt)
	})

	t.Run("should no-op when transaction is already deleted", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		tx.State = transaction.State_Deleted

		// act
		evt, err := tx.SetExpectedReimbursement(uuid.New(), values.Currency("USD"), decimal.NewFromInt(100))

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.ExpectedReimbursementSet{}, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())

		// act
		evt, err := tx.SetExpectedReimbursement(uuid.New(), values.Currency("USD"), decimal.NewFromInt(0))

		// assert
		require.ErrorIs(t, err, transaction.ErrNegativeOrNullAmount)
		assert.Equal(t, events.ExpectedReimbursementSet{}, evt)
	})
}

func TestRegisterExpense(t *testing.T) {
	t.Run("should emit money spent event when amount is positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		accountID := uuid.New()
		amount := decimal.NewFromInt(50)

		// act
		evt, err := tx.RegisterExpense(accountID, values.Currency("EUR"), amount, "food", "lunch")

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.MoneySpent{
			AccountID:   accountID,
			Currency:    values.Currency("EUR"),
			Amount:      amount,
			Category:    "food",
			Description: "lunch",
		}, evt)
	})

	t.Run("should no-op when transaction is already deleted", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		tx.State = transaction.State_Deleted

		// act
		evt, err := tx.RegisterExpense(uuid.New(), values.Currency("EUR"), decimal.NewFromInt(50), "food", "lunch")

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.MoneySpent{}, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())

		// act
		evt, err := tx.RegisterExpense(uuid.New(), values.Currency("EUR"), decimal.NewFromInt(-1), "food", "lunch")

		// assert
		require.ErrorIs(t, err, transaction.ErrNegativeOrNullAmount)
		assert.Equal(t, events.MoneySpent{}, evt)
	})
}

func TestRegisterIncome(t *testing.T) {
	t.Run("should emit money received event when amount is positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		accountID := uuid.New()
		amount := decimal.NewFromInt(80)

		// act
		evt, err := tx.RegisterIncome(accountID, values.Currency("GBP"), amount, "salary", "bonus")

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.MoneyReceived{
			AccountID:   accountID,
			Currency:    values.Currency("GBP"),
			Amount:      amount,
			Category:    "salary",
			Description: "bonus",
		}, evt)
	})

	t.Run("should no-op when transaction is already deleted", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		tx.State = transaction.State_Deleted

		// act
		evt, err := tx.RegisterIncome(uuid.New(), values.Currency("GBP"), decimal.NewFromInt(80), "salary", "bonus")

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.MoneyReceived{}, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())

		// act
		evt, err := tx.RegisterIncome(uuid.New(), values.Currency("GBP"), decimal.NewFromInt(0), "salary", "bonus")

		// assert
		require.ErrorIs(t, err, transaction.ErrNegativeOrNullAmount)
		assert.Equal(t, events.MoneyReceived{}, evt)
	})
}

func TestRegisterTransfer(t *testing.T) {
	t.Run("should emit money transferred event when payload is valid", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		fromAccountID := uuid.New()
		toAccountID := uuid.New()

		// act
		evt, err := tx.RegisterTransfer(
			fromAccountID,
			values.Currency("USD"),
			decimal.NewFromInt(100),
			toAccountID,
			values.Currency("EUR"),
			decimal.NewFromInt(90),
			"transfer",
			"currency exchange",
		)

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.MoneyTransfered{
			FromAccountID: fromAccountID,
			FromCurrency:  values.Currency("USD"),
			FromAmount:    decimal.NewFromInt(100),
			ToAccountID:   toAccountID,
			ToCurrency:    values.Currency("EUR"),
			ToAmount:      decimal.NewFromInt(90),
			Category:      "transfer",
			Description:   "currency exchange",
		}, evt)
	})

	t.Run("should no-op when transaction is already deleted", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		tx.State = transaction.State_Deleted

		// act
		evt, err := tx.RegisterTransfer(
			uuid.New(),
			values.Currency("USD"),
			decimal.NewFromInt(100),
			uuid.New(),
			values.Currency("EUR"),
			decimal.NewFromInt(90),
			"transfer",
			"currency exchange",
		)

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.MoneyTransfered{}, evt)
	})

	t.Run("should return error when any amount is not positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())

		// act
		evt, err := tx.RegisterTransfer(
			uuid.New(),
			values.Currency("USD"),
			decimal.NewFromInt(0),
			uuid.New(),
			values.Currency("EUR"),
			decimal.NewFromInt(90),
			"transfer",
			"currency exchange",
		)

		// assert
		require.ErrorIs(t, err, transaction.ErrNegativeOrNullAmount)
		assert.Equal(t, events.MoneyTransfered{}, evt)
	})

	t.Run("should return error when accounts are the same", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		accountID := uuid.New()

		// act
		evt, err := tx.RegisterTransfer(
			accountID,
			values.Currency("USD"),
			decimal.NewFromInt(100),
			accountID,
			values.Currency("EUR"),
			decimal.NewFromInt(90),
			"transfer",
			"currency exchange",
		)

		// assert
		require.ErrorIs(t, err, transaction.ErrInvalidAccount)
		assert.Equal(t, events.MoneyTransfered{}, evt)
	})

	t.Run("should return error when currencies match but amounts differ", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())

		// act
		evt, err := tx.RegisterTransfer(
			uuid.New(),
			values.Currency("USD"),
			decimal.NewFromInt(100),
			uuid.New(),
			values.Currency("USD"),
			decimal.NewFromInt(50),
			"transfer",
			"currency exchange",
		)

		// assert
		require.ErrorIs(t, err, transaction.ErrInvalidAmountOrCurrency)
		assert.Equal(t, events.MoneyTransfered{}, evt)
	})
}

func TestRegisterReimbursement(t *testing.T) {
	t.Run("should emit reimbursement received event when amount is positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		accountID := uuid.New()
		amount := decimal.NewFromInt(75)

		// act
		evt, err := tx.RegisterReimbursement(accountID, "company", values.Currency("USD"), amount)

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.ReimbursementReceived{
			AccountID: accountID,
			From:      "company",
			Currency:  values.Currency("USD"),
			Amount:    amount,
		}, evt)
	})

	t.Run("should no-op when transaction is already deleted", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())
		tx.State = transaction.State_Deleted

		// act
		evt, err := tx.RegisterReimbursement(uuid.New(), "company", values.Currency("USD"), decimal.NewFromInt(75))

		// assert
		require.NoError(t, err)
		assert.Equal(t, events.ReimbursementReceived{}, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		tx := transaction.New(uuid.New())

		// act
		evt, err := tx.RegisterReimbursement(uuid.New(), "company", values.Currency("USD"), decimal.NewFromInt(0))

		// assert
		require.ErrorIs(t, err, transaction.ErrNegativeOrNullAmount)
		assert.Equal(t, events.ReimbursementReceived{}, evt)
	})
}
