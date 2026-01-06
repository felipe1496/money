package constants

type TransactionType string

const (
	SimpleExpense TransactionType = "simple_expense"
	Income        TransactionType = "income"
	Installment   TransactionType = "installment"
)

type InstanceType string

const (
	ThisOne          = "one"
	ThisAndFollowing = "following"
	All              = "all"
)
