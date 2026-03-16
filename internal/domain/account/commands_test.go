package account_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/internal/domain/account"
	"github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/values"
)

func TestOpen(t *testing.T) {
	now := time.Now()
	t.Run("should emit opened event when account is unopened", func(t *testing.T) {
		// arrange
		id := uuid.New()
		acc := account.New(id)
		name := "Personal"
		currency := values.Currency("EUR")

		// act
		evt, err := acc.Open(name, currency, now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.Opened{
			AccountID:  id,
			Name:       name,
			Currency:   currency,
			HappenedAt: now,
		}, evt)
	})

	t.Run("should return error when account is already opened", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Opened

		// act
		evt, err := acc.Open("Personal", values.Currency("EUR"), now)

		// assert
		require.NoError(t, err)
		assert.Nil(t, evt)
	})
}

func TestUpdateName(t *testing.T) {
	now := time.Now()
	t.Run("should emit name updated event when account is opened", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Opened
		newName := "Main"

		// act
		evt, err := acc.UpdateName(newName, now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.NameUpdated{
			Name:       newName,
			HappenedAt: now,
		}, evt)
	})

	t.Run("should return error when account is not opened", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())

		// act
		evt, err := acc.UpdateName("Main", now)

		// assert
		require.ErrorIs(t, err, account.ErrAccountNotOpened)
		assert.Nil(t, evt)
	})
}

func TestDeposit(t *testing.T) {
	now := time.Now()
	t.Run("should emit money deposited event when payload is valid", func(t *testing.T) {
		// arrange
		id := uuid.New()
		acc := account.New(id)
		acc.State = account.State_Opened
		amount := decimal.NewFromInt(100)
		user := "user-123"

		// act
		evt, err := acc.Deposit(values.Currency("EUR"), amount, user, now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.MoneyDeposited{
			AccountID:  id,
			Currency:   values.Currency("EUR"),
			Amount:     amount,
			User:       user,
			HappenedAt: now,
		}, evt)
	})

	t.Run("should return error when account is not opened", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())

		// act
		evt, err := acc.Deposit(values.Currency("EUR"), decimal.NewFromInt(100), "user", now)

		// assert
		require.ErrorIs(t, err, account.ErrAccountNotOpened)
		assert.Nil(t, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Opened

		// act
		evt, err := acc.Deposit(values.Currency("EUR"), decimal.Zero, "user", now)

		// assert
		require.ErrorIs(t, err, account.ErrNegativeOrNullAmount)
		assert.Nil(t, evt)
	})
}

func TestWithdraw(t *testing.T) {
	now := time.Now()
	t.Run("should emit money withdrawn event when payload is valid", func(t *testing.T) {
		// arrange
		id := uuid.New()
		acc := account.New(id)
		acc.State = account.State_Opened
		amount := decimal.NewFromInt(50)
		user := "user-123"

		// act
		evt, err := acc.Withdraw(values.Currency("EUR"), amount, user, now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.MoneyWithdrawn{
			AccountID:  id,
			Currency:   values.Currency("EUR"),
			Amount:     amount,
			User:       user,
			HappenedAt: now,
		}, evt)
	})

	t.Run("should return error when account is not opened", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())

		// act
		evt, err := acc.Withdraw(values.Currency("EUR"), decimal.NewFromInt(50), "user", now)

		// assert
		require.ErrorIs(t, err, account.ErrAccountNotOpened)
		assert.Nil(t, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Opened

		// act
		evt, err := acc.Withdraw(values.Currency("EUR"), decimal.NewFromInt(-1), "user", now)

		// assert
		require.ErrorIs(t, err, account.ErrNegativeOrNullAmount)
		assert.Nil(t, evt)
	})
}
