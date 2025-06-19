package transaction

import (
	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type State int

const (
	State_Created State = iota
	State_Deleted
)

type Transaction struct {
	ID          uuid.UUID
	State       State
	Entries     []values.Entry
	Category    string
	Description string
}

func New() Transaction {
	return Transaction{
		ID:      uuid.Must(uuid.NewV7()),
		Entries: []values.Entry{},
	}
}
