package balances

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/somatom98/brokeli/internal/domain/account"
	account_events "github.com/somatom98/brokeli/internal/domain/account/events"
	"github.com/somatom98/brokeli/internal/domain/transaction"
	transaction_events "github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type BalancePeriod struct {
	Month    time.Time       `json:"month"`
	Currency values.Currency `json:"currency"`
	Amount   decimal.Decimal `json:"amount"`
}

type Repository interface {
	InsertBalanceUpdate(ctx context.Context, id uuid.UUID, accountID uuid.UUID, currency values.Currency, amount decimal.Decimal, userID string, valueDate time.Time) error
	GetBalancesByAccount(ctx context.Context, accountID uuid.UUID) ([]BalancePeriod, error)
	GetAllBalances(ctx context.Context) ([]BalancePeriod, error)
}

type Projection struct {
	repository Repository
}

func New(
	transactionES event_store.Store[*transaction.Transaction],
	accountES event_store.Store[*account.Account],
	repository Repository,
) *Projection {
	p := &Projection{
		repository: repository,
	}

	transactionES.Subscribe(context.Background(), p.HandleRecord)
	accountES.Subscribe(context.Background(), p.HandleRecord)

	return p
}

func (v *Projection) HandleRecord(ctx context.Context, record event_store.Record) error {
	var aggregateType string
	switch record.Type() {
	case transaction_events.TypeMoneySpent, transaction_events.TypeMoneyReceived, transaction_events.TypeReimbursementReceived:
		aggregateType = "Transaction"
	case account_events.TypeOpened, account_events.TypeMoneyDeposited, account_events.TypeMoneyWithdrawn:
		aggregateType = "Account"
	default:
		return nil
	}

	idStr := fmt.Sprintf("%s_%s_%d", aggregateType, record.AggregateID.String(), record.Version)
	id := uuid.NewMD5(uuid.NameSpaceOID, []byte(idStr))

	switch record.Type() {
	case transaction_events.TypeMoneySpent:
		return v.ApplyExpenseCreated(ctx, id, record.Content().(transaction_events.MoneySpent))
	case transaction_events.TypeMoneyReceived:
		return v.ApplyIncomeCreated(ctx, id, record.Content().(transaction_events.MoneyReceived))
	case transaction_events.TypeReimbursementReceived:
		return v.ApplyReimbursementReceived(ctx, id, record.Content().(transaction_events.ReimbursementReceived))
	case account_events.TypeMoneyDeposited:
		return v.ApplyMoneyDeposited(ctx, id, record.Content().(account_events.MoneyDeposited))
	case account_events.TypeMoneyWithdrawn:
		return v.ApplyMoneyWithdrawn(ctx, id, record.Content().(account_events.MoneyWithdrawn))
	}
	return nil
}

func (v *Projection) GetBalancesByAccount(ctx context.Context, accountID uuid.UUID) ([]BalancePeriod, error) {
	return v.repository.GetBalancesByAccount(ctx, accountID)
}

func (v *Projection) GetAllBalances(ctx context.Context) ([]BalancePeriod, error) {
	return v.repository.GetAllBalances(ctx)
}
