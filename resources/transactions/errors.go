package transactions

import (
	"fmt"
	"net/http"
	"rango-backend/utils"
)

var (
	ErrFailedToFetchEntries                 = fmt.Errorf("failed to fetch entries")
	ErrToCountEntries                       = fmt.Errorf("failed to count entries")
	ItWasNotPossibleDeleteTransactionErr    = utils.NewHTTPError(http.StatusInternalServerError, "It was not possible to delete transaction")
	TransactionNotFound                     = utils.NewHTTPError(http.StatusNotFound, "Transaction not found")
	AnErrorOccuredWhileFetchingTransactions = utils.NewHTTPError(http.StatusInternalServerError, "An error occured while fetching transactions")
)
