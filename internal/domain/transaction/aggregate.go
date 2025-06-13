package transaction

import (
	"github.com/google/uuid"
	"github.com/somatom98/brokeli/internal/domain/values"
)

type Transasction struct {
	ID          uuid.UUID
	Entries     []values.Entry
	Category    string
	Description string
}
