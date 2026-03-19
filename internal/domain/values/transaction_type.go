package values

type TransactionType string

const (
	TransactionType_Expense               TransactionType = "EXPENSE"
	TransactionType_Income                TransactionType = "INCOME"
	TransactionType_Transfer              TransactionType = "TRANSFER"
	TransactionType_Reimbursement         TransactionType = "REIMBURSEMENT"
	TransactionType_ExpectedReimbursement TransactionType = "EXPECTED_REIMBURSEMENT"
	TransactionType_Deposit               TransactionType = "DEPOSIT"
	TransactionType_Withdrawal            TransactionType = "WITHDRAWAL"
)
