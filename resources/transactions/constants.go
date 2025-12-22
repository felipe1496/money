package transactions

type TransactionType string

const (
	SimpleExpense TransactionType = "simple_expense"
	Income        TransactionType = "income"
	Installment   TransactionType = "installment"
)
