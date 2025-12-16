package transactions

import (
	"fmt"
	"rango-backend/utils"

	"github.com/Masterminds/squirrel"
	"github.com/oklog/ulid/v2"
)

type TransactionsRepo interface {
	CreateEntry(payload CreateEntryDTO, db utils.Executer) (Entry, error)
	CreateTransaction(payload CreateTransactionDTO, db utils.Executer) (Transaction, error)
	ListViewEntries(db utils.Executer, filter *utils.FilterBuilder) ([]ViewEntry, error)
	CountViewEntries(db utils.Executer, filter *utils.FilterBuilder) (int, error)
}

type TransactionsRepoImpl struct {
}

func NewTransactionsRepo(db utils.Executer) TransactionsRepo {
	return &TransactionsRepoImpl{}
}

func (r *TransactionsRepoImpl) CreateEntry(payload CreateEntryDTO, db utils.Executer) (Entry, error) {

	query, args, err := squirrel.Insert("entries").
		Columns("id", "transaction_id", "amount", "period").
		Values(ulid.Make().String(), payload.TransactionID, payload.Amount, payload.Period).
		Suffix("RETURNING id, transaction_id, amount, period, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return Entry{}, err
	}

	var entry Entry
	err = db.QueryRow(query, args...).Scan(
		&entry.ID,
		&entry.TransactionID,
		&entry.Amount,
		&entry.Period,
		&entry.CreatedAt,
	)
	return entry, err
}

func (r *TransactionsRepoImpl) CreateTransaction(payload CreateTransactionDTO, db utils.Executer) (Transaction, error) {

	query, args, err := squirrel.Insert("transactions").
		Columns("id", "user_id", "category", "name", "description").
		Values(ulid.Make().String(), payload.UserID, payload.Type, payload.Name, &payload.Description).
		Suffix("RETURNING id, user_id, category, name, description, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return Transaction{}, err
	}

	var transaction Transaction
	err = db.QueryRow(query, args...).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Type,
		&transaction.Name,
		&transaction.Description,
		&transaction.CreatedAt,
	)
	return transaction, err
}

func (r *TransactionsRepoImpl) ListViewEntries(db utils.Executer, filter *utils.FilterBuilder) ([]ViewEntry, error) {
	query := squirrel.Select("id", "transaction_id", "name", "description", "amount", "period", "user_id", "category", "total_amount", "installment", "total_installments").
		From("v_entries").
		PlaceholderFormat(squirrel.Dollar)

	query, err := utils.ApplyFilterToSquirrel(query, filter)
	if err != nil {
		return nil, err
	}

	sql, args, err := query.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := db.Query(sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var entries []ViewEntry = []ViewEntry{}
	for rows.Next() {
		var entry ViewEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.TransactionID,
			&entry.Name,
			&entry.Description,
			&entry.Amount,
			&entry.Period,
			&entry.UserID,
			&entry.Type,
			&entry.TotalAmount,
			&entry.Installment,
			&entry.TotalInstallments,
		); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *TransactionsRepoImpl) CountViewEntries(db utils.Executer, filter *utils.FilterBuilder) (int, error) {
	countQuery := squirrel.
		Select("COUNT(*)").
		From("v_entries").
		PlaceholderFormat(squirrel.Dollar)

	countQuery, err := utils.ApplyFilterToSquirrel(
		countQuery,
		filter,
	)

	if err != nil {
		return 0, err
	}

	sql, args, err := countQuery.ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	err = db.QueryRow(sql, args...).Scan(&count)
	fmt.Println(args)
	if err != nil {
		return 0, err
	}

	return count, nil
}
