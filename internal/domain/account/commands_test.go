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

func TestAccountCreate(t *testing.T) {
	t.Run("should emit created event when account is new", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		createdAt := time.Date(2024, 10, 19, 12, 0, 0, 0, time.UTC)

		// act
		evt, err := acc.Create(createdAt)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.Created{Time: createdAt}, evt)
	})

	t.Run("should not emit event when account is already created", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Created

		// act
		evt, err := acc.Create(time.Date(2024, 10, 19, 12, 0, 0, 0, time.UTC))

		// assert
		require.NoError(t, err)
		assert.Nil(t, evt)
	})

	t.Run("should not return validation errors when timestamp is zero", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())

		// act
		evt, err := acc.Create(time.Time{})

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.Created{Time: time.Time{}}, evt)
	})
}

func TestAccountDeposit(t *testing.T) {
	t.Run("should emit deposit event when amount is positive", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		amount := decimal.NewFromInt(50)
		now := time.Date(2024, 11, 1, 9, 30, 0, 0, time.UTC)

		// act
		evt, err := acc.Deposit("user-123", values.Currency("USD"), amount, now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.MoneyDeposited{
			User:     "user-123",
			Currency: values.Currency("USD"),
			Amount:   amount,
			Time:     now,
		}, evt)
	})

	t.Run("should no-op when account is already closed", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Closed

		// act
		evt, err := acc.Deposit("user-123", values.Currency("USD"), decimal.NewFromInt(50), time.Date(2024, 11, 1, 9, 30, 0, 0, time.UTC))

		// assert
		require.NoError(t, err)
		assert.Nil(t, evt)
	})

	t.Run("should return error when amount is not positive", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())

		// act
		evt, err := acc.Deposit("user-123", values.Currency("USD"), decimal.NewFromInt(0), time.Date(2024, 11, 1, 9, 30, 0, 0, time.UTC))

		// assert
		require.ErrorIs(t, err, account.ErrNegativeOrNullAmount)
		assert.Nil(t, evt)
	})
}

func TestAccountClose(t *testing.T) {
	t.Run("should emit account closed event when account is open", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Created
		now := time.Date(2024, 12, 10, 8, 0, 0, 0, time.UTC)

		// act
		evt, err := acc.Close(now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.AccountClosed{Time: now}, evt)
	})

	t.Run("should no-op when account is already closed", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		acc.State = account.State_Closed

		// act
		evt, err := acc.Close(time.Date(2024, 12, 10, 8, 0, 0, 0, time.UTC))

		// assert
		require.NoError(t, err)
		assert.Nil(t, evt)
	})

	t.Run("should not return validation errors when closing unknown account", func(t *testing.T) {
		// arrange
		acc := account.New(uuid.New())
		now := time.Date(2024, 12, 10, 8, 0, 0, 0, time.UTC)

		// act
		evt, err := acc.Close(now)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &events.AccountClosed{Time: now}, evt)
	})
}
