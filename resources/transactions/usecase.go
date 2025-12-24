package transactions

import (
	"database/sql"
	"net/http"
	"rango-backend/utils"
)

type TransactionsUseCase interface {
	CreateSimpleExpense(payload CreateSimpleExpenseDTO) (string, error)
	ListViewEntries(filter *utils.QueryOptsBuilder) ([]ViewEntry, error)
	CountViewEntries(filter *utils.QueryOptsBuilder) (int, error)
	DeleteTransactionById(id string) error
	CreateIncome(payload CreateIncomeDTO) (string, error)
	CreateInstallment(payload CreateInstallmentDTO) ([]ViewEntry, error)
}

type TransactionsUseCaseImpl struct {
	repo TransactionsRepo
	db   *sql.DB
}

func NewTransactionsUseCase(repo TransactionsRepo, db *sql.DB) TransactionsUseCase {
	return &TransactionsUseCaseImpl{
		repo: repo,
		db:   db,
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

	transaction, err := uc.repo.CreateTransaction(CreateTransactionDTO{
		UserID:      payload.UserID,
		Type:        SimpleExpense,
		Name:        payload.Name,
		Description: &payload.Description,
	}, conn)

	if err != nil {
		return "", err
	}

	if payload.Amount > 0 {
		payload.Amount = payload.Amount * -1
	}

	_, err = uc.repo.CreateEntry(CreateEntryDTO{
		TransactionID: transaction.ID,
		Amount:        payload.Amount,
		ReferenceDate: payload.ReferenceDate,
	}, conn)

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

func (uc *TransactionsUseCaseImpl) CreateIncome(payload CreateIncomeDTO) (string, error) {
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

	transaction, err := uc.repo.CreateTransaction(CreateTransactionDTO{
		UserID:      payload.UserID,
		Type:        Income,
		Name:        payload.Name,
		Description: &payload.Description,
	}, conn)

	if err != nil {
		return "", err
	}

	_, err = uc.repo.CreateEntry(CreateEntryDTO{
		TransactionID: transaction.ID,
		Amount:        payload.Amount,
		ReferenceDate: payload.ReferenceDate,
	}, conn)

	if err != nil {
		return "", err
	}

	if err := conn.Commit(); err != nil {
		return "", err
	}

	return transaction.ID, nil
}

func (uc *TransactionsUseCaseImpl) CreateInstallment(payload CreateInstallmentDTO) ([]ViewEntry, error) {
	conn, err := uc.db.Begin()

	if err != nil {
		return []ViewEntry{}, utils.NewHTTPError(http.StatusInternalServerError, "An error occurred while trying to start a connection with the database")
	}

	defer func() {
		if err != nil {
			conn.Rollback()
			return
		}
	}()

	transaction, err := uc.repo.CreateTransaction(CreateTransactionDTO{
		UserID:      payload.UserID,
		Type:        Installment,
		Name:        payload.Name,
		Description: &payload.Description,
	}, conn)

	if err != nil {
		return []ViewEntry{}, utils.NewHTTPError(http.StatusInternalServerError, "An error occurred while trying to create the transaction")
	}
	totalAmountCents := int(payload.TotalAmount * 100)

	baseTotalAmount := totalAmountCents / payload.TotalInstallments
	rest := totalAmountCents % payload.TotalInstallments

	for i := 0; i < payload.TotalInstallments; i++ {
		nextRefDate := payload.ReferenceDate.AddDate(0, i, 0)
		amount := baseTotalAmount
		if i == 0 {
			amount += rest
		}
		entryDTO := CreateEntryDTO{
			TransactionID: transaction.ID,
			Amount:        float64(amount) / 100,
			ReferenceDate: nextRefDate,
		}

		_, err = uc.repo.CreateEntry(entryDTO, conn)
	}

	if err != nil {
		return []ViewEntry{}, err
	}

	if err := conn.Commit(); err != nil {
		return []ViewEntry{}, err
	}

	entries, _ := uc.ListViewEntries(utils.QueryOpts().And("transaction_id", "eq", transaction.ID))

	return entries, nil
}
