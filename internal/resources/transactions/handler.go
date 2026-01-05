package transactions

import (
	"database/sql"
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/gin-gonic/gin"
)

type API struct {
	transactionsUseCase TransactionsUseCase
}

func NewHandler(db *sql.DB) *API {
	return &API{
		transactionsUseCase: NewTransactionsUseCase(NewTransactionsRepo(db), db),
	}
}

// @Summary Create a Simple Expense
// @Description Create a Simple Expense transaction and entry
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateSimpleExpenseRequest true "Simple Expense payload"
// @Success 201 {object} CreateSimpleExpenseResponse "Simple Expense created"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/simple-expense [post]
func (api *API) CreateSimpleExpense(ctx *gin.Context) {
	var body CreateSimpleExpenseRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	id, err := api.transactionsUseCase.CreateSimpleExpense(CreateSimpleExpenseDTO{
		Name:          body.Name,
		Amount:        body.Amount,
		ReferenceDate: body.ReferenceDate,
		Description:   body.Description,
		UserID:        ctx.GetString("user_id"),
		CategoryID:    body.CategoryID,
	})

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	entries, err := api.transactionsUseCase.ListViewEntries(utils.QueryOpts().And("transaction_id", "eq", id))

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, CreateSimpleExpenseResponse{
		Data: CreateSimpleExpenseResponseData{
			Entry: entries[0],
		},
	})
}

// @Summary Create an income
// @Description Create an income transactions and entry
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateIncomeRequest true "Income payload"
// @Success 201 {object} CreateIncomeResponse "Income transaction created"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/income [post]
func (api *API) CreateIncome(ctx *gin.Context) {
	var body CreateIncomeRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	id, err := api.transactionsUseCase.CreateIncome(CreateIncomeDTO{
		Name:          body.Name,
		Amount:        body.Amount,
		ReferenceDate: body.ReferenceDate,
		Description:   body.Description,
		UserID:        ctx.GetString("user_id"),
		CategoryID:    body.CategoryID,
	})

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	entries, err := api.transactionsUseCase.ListViewEntries(utils.QueryOpts().And("transaction_id", "eq", id))

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, CreateIncomeResponse{
		Data: CreateIncomeResponseData{
			Entry: entries[0],
		},
	})
}

// @Summary List entries
// @Description List a detailed view of entries joined with transactions for a given period
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "period in format YYYYMM"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param order_by query string false "Sort field" example(name)
// @Param order query string false "Sort order (asc/desc)" Enums(asc, desc) default(asc)
// @Success 200 {object} ListEntriesResponse "List of entries"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/entries/{period} [get]
func (api *API) ListViewEntries(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	period := ctx.Param("period")
	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	queryOpts := ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).And("user_id", "eq", userID).
		And("period", "eq", period)

	entries, err := api.transactionsUseCase.ListViewEntries(queryOpts)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.transactionsUseCase.CountViewEntries(utils.QueryOpts().
		And("user_id", "eq", userID).
		And("period", "eq", period))

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	nextPage := len(entries) > perPage

	if nextPage {
		entries = entries[:len(entries)-1]
	}

	totalPages := (count + perPage - 1) / perPage

	ctx.JSON(http.StatusOK, ListEntriesResponse{
		Data: ListEntriesResponseData{
			Entries: entries,
		},
		Query: utils.QueryMeta{
			Page:       page,
			PerPage:    perPage,
			NextPage:   nextPage,
			TotalPages: totalPages,
			TotalItems: count,
		},
	})
}

// @Summary Delete Transaction By ID
// @Description Delete a transaction and all entries related by the ID of the transaction
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Success 204 "Transaction deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/{transaction_id} [delete]
func (api *API) DeleteTransaction(ctx *gin.Context) {
	id := ctx.Param("transaction_id")

	err := api.transactionsUseCase.DeleteTransactionById(id)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary Create an installment
// @Description Create an installment transactions and entries
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateInstallmentRequest true "Installment payload"
// @Success 201 {object} CreateInstallmentResponse "Transaction created"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/installment [post]
func (api *API) CreateInstallment(ctx *gin.Context) {
	var body CreateInstallmentRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	entries, err := api.transactionsUseCase.CreateInstallment(CreateInstallmentDTO{
		Name:              body.Name,
		TotalAmount:       body.TotalAmount,
		TotalInstallments: body.TotalInstallments,
		ReferenceDate:     body.ReferenceDate,
		Description:       body.Description,
		UserID:            ctx.GetString("user_id"),
		CategoryID:        body.CategoryID,
	})

	ctx.JSON(http.StatusCreated, CreateInstallmentResponse{
		Data: CreateInstallmentResponseData{
			Entries: entries,
		},
	})
}

// @Summary Update a simple expense
// @Description Update a simple expense
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction ID"
// @Param body body UpdateSimpleExpenseRequest true "Simple expense payload"
// @Success 200 {object} UpdateSimpleExpenseResponse "Simple expense updated"
// @Failure 400 {object} utils.HTTPError "Bad request"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /transactions/simple-expense/{transaction_id} [patch]
func (api *API) UpdateSimpleExpense(ctx *gin.Context) {
	transactionID := ctx.Param("transaction_id")

	var body UpdateSimpleExpenseRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	entry, err := api.transactionsUseCase.UpdateSimpleExpense(transactionID, UpdateSimpleExpenseDTO{
		Name:          body.Name,
		Description:   body.Description,
		Amount:        body.Amount,
		ReferenceDate: body.ReferenceDate,
		CategoryID:    body.CategoryID,
	})

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, UpdateSimpleExpenseResponse{
		Data: UpdateSimpleExpenseResponseData{
			Entry: entry,
		},
	})
}
