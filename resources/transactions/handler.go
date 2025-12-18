package transactions

import (
	"database/sql"
	"net/http"
	"rango-backend/utils"

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

func (api *API) CreateSimpleExpense(ctx *gin.Context) {
	var body CreateSimpleExpenseRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	id, err := api.transactionsUseCase.CreateSimpleTransaction(CreateSimpleExpenseDTO{
		Name:        body.Name,
		Amount:      body.Amount,
		Period:      body.Period,
		Description: body.Description,
		UserID:      ctx.GetString("user_id"),
	})

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	entries, err := api.transactionsUseCase.ListViewEntries(utils.CreateFilter().And("transaction_id", "eq", id))

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"data": gin.H{
			"entry": entries[0],
		},
	})
}

func (api *API) ListViewEntries(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	period := ctx.Param("period")
	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	filter := ctx.MustGet("filter").(*utils.FilterBuilder).And("user_id", "eq", userID).
		And("period", "eq", period).OrderBy("created_at", "desc")

	entries, err := api.transactionsUseCase.ListViewEntries(filter)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.transactionsUseCase.CountViewEntries(utils.CreateFilter().
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

	ctx.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"entries": entries,
		},
		"query": gin.H{
			"page":        page,
			"per_page":    perPage,
			"next_page":   nextPage,
			"total_pages": totalPages,
			"total_items": count,
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
