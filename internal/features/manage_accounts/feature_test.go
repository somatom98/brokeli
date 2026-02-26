package manage_accounts_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/features/manage_accounts"
	"github.com/somatom98/brokeli/pkg/event_store"
	"github.com/stretchr/testify/assert"
)

type DispatcherMock struct {
	Withdrawals []withdrawalCall
	Deposits    []depositCall
}

type withdrawalCall struct {
	ID       uuid.UUID
	Currency values.Currency
	Amount   decimal.Decimal
	User     string
}

type depositCall struct {
	ID       uuid.UUID
	Currency values.Currency
	Amount   decimal.Decimal
	User     string
}

func (m *DispatcherMock) Open(ctx context.Context, id uuid.UUID, name string, currency values.Currency) error { return nil }
func (m *DispatcherMock) UpdateName(ctx context.Context, id uuid.UUID, name string) error { return nil }
func (m *DispatcherMock) Deposit(ctx context.Context, id uuid.UUID, currency values.Currency, amount decimal.Decimal, user string) error {
	m.Deposits = append(m.Deposits, depositCall{
		ID:       id,
		Currency: currency,
		Amount:   amount,
		User:     user,
	})
	return nil
}
func (m *DispatcherMock) Withdraw(ctx context.Context, id uuid.UUID, currency values.Currency, amount decimal.Decimal, user string) error {
	m.Withdrawals = append(m.Withdrawals, withdrawalCall{
		ID:       id,
		Currency: currency,
		Amount:   amount,
		User:     user,
	})
	return nil
}
func (m *DispatcherMock) Transfer(ctx context.Context, fromID uuid.UUID, toID uuid.UUID, fromCurrency values.Currency, toCurrency values.Currency, fromAmount decimal.Decimal, toAmount decimal.Decimal, user string) error {
	return nil
}

func TestManageAccounts_TransferEventHandler(t *testing.T) {
	// arrange
	transactionES := event_store.NewInMemory[*transaction.Transaction](transaction.New)
	dispatcher := &DispatcherMock{}
	feature := manage_accounts.New(&http.ServeMux{}, nil, dispatcher, transactionES)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	feature.Setup(ctx)
	
	fromID := uuid.New()
	toID := uuid.New()
	amount := decimal.NewFromInt(100)
	
	event := transaction_events.MoneyTransfered{
		FromAccountID: fromID,
		ToAccountID:   toID,
		FromCurrency:  values.Currency("EUR"),
		ToCurrency:    values.Currency("EUR"),
		FromAmount:    amount,
		ToAmount:      amount,
	}
	
	// act
	err := transactionES.Append(ctx, event_store.Record{
		AggregateID: uuid.New(),
		Version:     1,
		Event:       event,
	})
	assert.NoError(t, err)
	
	// assert
	assert.Eventually(t, func() bool {
		return len(dispatcher.Withdrawals) == 1 && len(dispatcher.Deposits) == 1
	}, 1*time.Second, 10*time.Millisecond)
	
	withdrawal := dispatcher.Withdrawals[0]
	assert.Equal(t, fromID, withdrawal.ID)
	assert.Equal(t, amount, withdrawal.Amount)
	assert.Equal(t, "system", withdrawal.User)

	deposit := dispatcher.Deposits[0]
	assert.Equal(t, toID, deposit.ID)
	assert.Equal(t, amount, deposit.Amount)
	assert.Equal(t, "system", deposit.User)
}
