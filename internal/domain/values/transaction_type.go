package values

type TransactionType int

const (
	TransactionType_Expense TransactionType = iota
	TransactionType_Income
	TransactionType_Transfer
	TransactionType_Reimbursement
	TransactionType_ExpectedReimbursement
)
