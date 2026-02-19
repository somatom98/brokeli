package import_transactions_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/somatom98/brokeli/internal/domain/transaction"
	"github.com/somatom98/brokeli/internal/domain/transaction/events"
	"github.com/somatom98/brokeli/internal/domain/values"
	"github.com/somatom98/brokeli/internal/features/import_transactions"
	"github.com/somatom98/brokeli/pkg/event_store"
)

type SpyStore struct {
	records []event_store.Record
}

func (s *SpyStore) Append(ctx context.Context, record event_store.Record) error {
	s.records = append(s.records, record)
	return nil
}

func (s *SpyStore) GetAggregate(ctx context.Context, id uuid.UUID) (*transaction.Transaction, error) {
	return transaction.New(id), nil
}

func (s *SpyStore) Subscribe(ctx context.Context) <-chan event_store.Record {
	return nil
}

func TestImportTransactions(t *testing.T) {
	// arrange
	spy := &SpyStore{}
	dispatcher := transaction.NewDispatcher(spy)
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

	// act
	err = feature.ImportTransactions(context.Background(), tmpFile.Name())

	// assert
	require.NoError(t, err)
	assert.Len(t, spy.records, 4)

	// Expected IDs
	lunarID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Lunar"))
	intesaID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Intesa San Paolo"))
	directaID := uuid.NewMD5(uuid.NameSpaceOID, []byte("Directa"))

	// 1. Expense
	assert.IsType(t, &events.MoneySpent{}, spy.records[0].Event)
	spent := spy.records[0].Event.(*events.MoneySpent)
	assert.Equal(t, lunarID, spent.AccountID)
	assert.Equal(t, values.Currency("DKK"), spent.Currency)
	assert.True(t, decimal.NewFromFloat(148).Equal(spent.Amount))

	// 2. Income
	assert.IsType(t, &events.MoneyReceived{}, spy.records[1].Event)
	received := spy.records[1].Event.(*events.MoneyReceived)
	assert.Equal(t, lunarID, received.AccountID)
	assert.Equal(t, values.Currency("DKK"), received.Currency)
	assert.True(t, decimal.NewFromFloat(123.82).Equal(received.Amount))

	// 3. Transfer
	assert.IsType(t, &events.MoneyTransfered{}, spy.records[2].Event)
	transferred := spy.records[2].Event.(*events.MoneyTransfered)
	assert.Equal(t, intesaID, transferred.FromAccountID)
	assert.Equal(t, directaID, transferred.ToAccountID)
	assert.Equal(t, values.Currency("EUR"), transferred.FromCurrency)
	assert.Equal(t, values.Currency("EUR"), transferred.ToCurrency)

	// 4. Reimbursement
	assert.IsType(t, &events.ReimbursementReceived{}, spy.records[3].Event)
	reimbursement := spy.records[3].Event.(*events.ReimbursementReceived)
	assert.Equal(t, lunarID, reimbursement.AccountID)
	assert.Equal(t, "Reimbursement dirt", reimbursement.From)
}

func TestImportTransactions_RealFile(t *testing.T) {
	// arrange
	spy := &SpyStore{}
	dispatcher := transaction.NewDispatcher(spy)
	feature := import_transactions.New(nil, dispatcher)

	// act
	err := feature.ImportTransactions(context.Background(), "transactions.csv")

	// assert
	require.NoError(t, err)
	assert.Greater(t, len(spy.records), 0)
	t.Logf("Imported %d transactions successfully", len(spy.records))
}
