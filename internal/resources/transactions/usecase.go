package transactions

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/felipe1496/open-wallet/internal/constants"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/utils"
)

type TransactionsUseCase interface {
	ListViewEntries(filter *utils.QueryOptsBuilder) ([]ViewEntry, error)
	CountViewEntries(filter *utils.QueryOptsBuilder) (int, error)
	DeleteTransactionById(id string) error
	CreateTransaction(payload CreateTransactionDTO2) (Transaction, error)
	UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO2) (Transaction, error)
}

type TransactionsUseCaseImpl struct {
	repo           TransactionsRepo
	categoriesRepo categories.CategoriesRepo
	db             *sql.DB
}

func NewTransactionsUseCase(repo TransactionsRepo, categoriesRepo categories.CategoriesRepo, db *sql.DB) TransactionsUseCase {
	return &TransactionsUseCaseImpl{
		repo,
		categoriesRepo,
		db,
	}
}

func (uc *TransactionsUseCaseImpl) CreateSimpleExpense(payload CreateSimpleExpenseDTO) (string, error) {
	conn, err := uc.db.Begin()

	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}
	}()

	transaction, err := uc.repo.CreateTransaction(conn, CreateTransactionDTO{
		UserID:      payload.UserID,
		Type:        constants.SimpleExpense,
		Name:        payload.Name,
		Description: &payload.Description,
		CategoryID:  payload.CategoryID,
	})

	if err != nil {
		return "", err
	}

	if payload.Amount > 0 {
		payload.Amount = payload.Amount * -1
	}

	_, err = uc.repo.CreateEntry(conn, CreateEntryDTO{
		TransactionID: transaction.ID,
		Amount:        payload.Amount,
		ReferenceDate: payload.ReferenceDate,
	})

	if err != nil {
		return "", err
	}

	if err := conn.Commit(); err != nil {
		return "", err
	}

	return transaction.ID, nil
}

func (uc *TransactionsUseCaseImpl) ListViewEntries(filter *utils.QueryOptsBuilder) ([]ViewEntry, error) {
	entries, err := uc.repo.ListViewEntries(uc.db, filter)

	if err != nil {
		return []ViewEntry{}, ErrFailedToFetchEntries
	}

	return entries, nil
}

func (uc *TransactionsUseCaseImpl) CountViewEntries(filter *utils.QueryOptsBuilder) (int, error) {
	count, err := uc.repo.CountViewEntries(uc.db, filter)

	if err != nil {
		return 0, ErrToCountEntries
	}

	return count, nil
}

func (uc *TransactionsUseCaseImpl) DeleteTransactionById(id string) error {
	transactionExists, err := uc.repo.ListTransactions(uc.db, utils.QueryOpts().And("id", "eq", id))

	if err != nil {
		return AnErrorOccuredWhileFetchingTransactions
	}

	if len(transactionExists) == 0 {
		return TransactionNotFound
	}

	err = uc.repo.DeleteTransactionById(uc.db, id)

	if err != nil {
		return ItWasNotPossibleDeleteTransactionErr
	}

	return nil
}

type validadeTransactionProps struct {
	Entries []validateTransactionPropsEntry
	Type    constants.TransactionType
}

type validateTransactionPropsEntry struct {
	Amount        float64
	ReferenceDate string
}

func validateTransaction(payload validadeTransactionProps) error {
	for i, refEntry := range payload.Entries {
		iRefDate, _ := time.Parse("2006-01-02", refEntry.ReferenceDate)
		iPeriod := iRefDate.Format("200601")
		for j, currEntry := range payload.Entries {
			if i != j {
				jRefDate, _ := time.Parse("2006-01-02", currEntry.ReferenceDate)
				jPeriod := jRefDate.Format("200601")
				if iPeriod == jPeriod {

					return utils.NewHTTPError(http.StatusBadRequest, "entries must be in different periods")
				}
			}
		}
	}

	switch payload.Type {
	case constants.SimpleExpense:
		{
			if len(payload.Entries) > 1 {
				return utils.NewHTTPError(http.StatusBadRequest, "expense must have only one entry")
			}
		}
	case constants.Income:
		{
			if len(payload.Entries) > 1 {
				return utils.NewHTTPError(http.StatusBadRequest, "income must have only one entry")
			}
		}
	case constants.Installment:
		{
			if len(payload.Entries) < 2 {
				return utils.NewHTTPError(http.StatusBadRequest, "installment must have at least two entries")
			}
		}
	}
	return nil
}

func (uc *TransactionsUseCaseImpl) CreateTransaction(payload CreateTransactionDTO2) (Transaction, error) {

	err := validateTransaction(validadeTransactionProps{
		Entries: func() []validateTransactionPropsEntry {
			entries := make([]validateTransactionPropsEntry, len(payload.Entries))
			for i, entry := range payload.Entries {
				entries[i] = validateTransactionPropsEntry{
					Amount:        entry.Amount,
					ReferenceDate: entry.ReferenceDate,
				}
			}
			return entries
		}(),
		Type: payload.Type,
	})
	if err != nil {
		return Transaction{}, err
	}

	tx, err := uc.db.Begin()

	if err != nil {
		return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to begin transaction")
	}

	transaction, err := uc.repo.CreateTransaction(tx, CreateTransactionDTO{
		UserID:      payload.UserID,
		Type:        payload.Type,
		Name:        payload.Name,
		Description: payload.Note,
		CategoryID:  payload.CategoryID,
	})

	if err != nil {
		tx.Rollback()
		return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create transaction")
	}

	for _, entry := range payload.Entries {
		if (payload.Type == constants.SimpleExpense || payload.Type == constants.Installment) && entry.Amount > 0 {
			entry.Amount = entry.Amount * -1
		} else if payload.Type == constants.Income && entry.Amount < 0 {
			entry.Amount = entry.Amount * -1
		}
		_, err = uc.repo.CreateEntry(tx, CreateEntryDTO{
			TransactionID: transaction.ID,
			Amount:        entry.Amount,
			ReferenceDate: entry.ReferenceDate,
		})

		if err != nil {
			tx.Rollback()
			return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
		}
	}

	err = tx.Commit()

	if err != nil {
		return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to commit transaction")
	}

	return transaction, nil
}

func (uc *TransactionsUseCaseImpl) UpdateTransaction(transactionID string, userID string, payload UpdateTransactionDTO2) (t Transaction, err error) {
	tx, err := uc.db.Begin()
	defer func() {
		if tx == nil {
			return
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if err != nil {
		return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to start transaction")
	}

	exists, err := uc.repo.ListViewEntries(tx, utils.QueryOpts().
		And("transaction_id", "eq", transactionID).
		And("user_id", "eq", userID))
	if err != nil {
		return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if transaction exists")
	}

	if len(exists) == 0 {
		return Transaction{}, utils.NewHTTPError(http.StatusNotFound, "transaction not found")
	}

	fmt.Println(">>> UpdateTransaction - update: ", payload.Update)
	if payload.CategoryID != nil && utils.Contains(payload.Update, "category_id") {
		categoryExists, err := uc.categoriesRepo.List(tx, utils.QueryOpts().
			And("id", "eq", *payload.CategoryID).
			And("user_id", "eq", userID))
		if err != nil {
			return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to check if category exists")
		}

		if len(categoryExists) == 0 {
			return Transaction{}, utils.NewHTTPError(http.StatusNotFound, "category not found")
		}
	}

	if utils.ContainsSome(payload.Update, []string{"name", "note", "category_id"}) {
		_, err = uc.repo.UpdateTransaction(tx, transactionID, UpdateTransactionDTO2{
			Update:     payload.Update,
			Name:       payload.Name,
			Note:       payload.Note,
			CategoryID: payload.CategoryID,
		})
		if err != nil {
			return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to update transaction")
		}
	}

	if payload.Entries != nil && utils.Contains(payload.Update, "entries") {
		err = validateTransaction(validadeTransactionProps{
			Entries: func() []validateTransactionPropsEntry {
				entries := make([]validateTransactionPropsEntry, 0)
				if payload.Entries != nil {
					for _, entry := range *payload.Entries {
						entries = append(entries, validateTransactionPropsEntry{
							Amount:        entry.Amount,
							ReferenceDate: entry.ReferenceDate,
						})
					}
				}
				return entries
			}(),
			Type: exists[0].Type,
		})

		if err != nil {
			return Transaction{}, err
		}

		err = uc.repo.DeleteEntry(tx, utils.QueryOpts().
			And("transaction_id", "eq", exists[0].TransactionID))
		if err != nil {
			return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to delete previous entries")
		}

		for i, entry := range *payload.Entries {
			fmt.Printf(">>> Calling CreateEntry [%d]\n", i)
			_, err = uc.repo.CreateEntry(tx, CreateEntryDTO{
				TransactionID: exists[0].TransactionID,
				Amount:        entry.Amount,
				ReferenceDate: entry.ReferenceDate,
			})
			if err != nil {
				return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to create entry")
			}
		}
	}

	transactions, err := uc.repo.ListTransactions(tx, utils.QueryOpts().
		And("id", "eq", exists[0].TransactionID))
	if err != nil {
		return Transaction{}, utils.NewHTTPError(http.StatusInternalServerError, "failed to list transaction")
	}

	return transactions[0], nil
}
