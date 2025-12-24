package categories

import (
	"database/sql"
	"net/http"
	"rango-backend/utils"

	"github.com/gin-gonic/gin"
)

type API struct {
	categoriesUseCase CategoriesUseCase
}

func NewHandler(db *sql.DB) *API {
	return &API{
		categoriesUseCase: NewCategoriesUseCase(NewCategoriesRepo(db), db),
	}
}

// @Summary Create a category
// @Description Create a category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateCategoryRequest true "Category payload"
// @Success 201 {object} CreateCategoryResponse "Category created"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories [post]
func (api *API) Create(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	var body CreateCategoryRequest

	err := ctx.ShouldBindJSON(&body)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusBadRequest, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	category, err := api.categoriesUseCase.Create(CreateCategoryDTO{
		UserID: userID,
		Name:   body.Name,
		Color:  body.Color,
	})

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusCreated, CreateCategoryResponse{
		Data: CreateCategoryResponseData{
			Category: category,
		},
	})
}

// @Summary List categories
// @Description List categories
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param sort query string false "Sort field" example(name)
// @Param order query string false "Sort order (asc/desc)" Enums(asc, desc) default(asc)
// @Success 200 {object} ListCategoriesResponse "List of categories"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories [get]
func (api *API) List(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	queryOpts := ctx.MustGet("query_opts").(*utils.QueryOptsBuilder).And("user_id", "eq", userID)

	categories, err := api.categoriesUseCase.List(queryOpts)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	count, err := api.categoriesUseCase.Count(queryOpts)
	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	page := ctx.GetInt("page")
	perPage := ctx.GetInt("per_page")
	nextPage := len(categories) > perPage
	totalPages := (count + perPage - 1) / perPage

	if nextPage {
		categories = categories[:len(categories)-1]
	}

	ctx.JSON(http.StatusOK, ListCategoriesResponse{
		Data: ListCategoriesResponseData{
			Categories: categories,
		},
		Query: utils.QueryMeta{
			NextPage:   nextPage,
			Page:       page,
			PerPage:    perPage,
			TotalItems: count,
			TotalPages: totalPages,
		},
	})
}

// @Summary Delete Category By ID
// @Description Delete a category
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category_id path string true "category ID"
// @Success 204 "Category deleted"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 404 {object} utils.HTTPError "Not found"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories/{category_id} [delete]
func (api *API) DeleteByID(ctx *gin.Context) {
	id := ctx.Param("category_id")

	err := api.categoriesUseCase.DeleteByID(id)

	if err != nil {
		apiErr := err.(*utils.HTTPError)
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// @Summary List categories with amount per period
// @Description List categories with amount per period
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param period path string true "period"
// @Success 200 {object} ListCategoryAmountPerPeriodResponse "List of categories with amount per period"
// @Failure 401 {object} utils.HTTPError "Unauthorized"
// @Failure 500 {object} utils.HTTPError "Internal server error"
// @Router /categories/{period} [get]
func (api *API) ListCategoryAmountPerPeriod(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	period := ctx.Param("period")
	queryOpts := utils.QueryOpts().And("user_id", "eq", userID).And("period", "eq", period)

	categories, err := api.categoriesUseCase.ListCategoryAmountPerPeriod(queryOpts)

	if err != nil {
		apiErr := utils.NewHTTPError(http.StatusInternalServerError, err.Error())
		ctx.JSON(apiErr.StatusCode, apiErr)
		return
	}

	ctx.JSON(http.StatusOK, ListCategoryAmountPerPeriodResponse{
		Data: ListCategoryAmountPerPeriodResponseData{
			Categories: categories,
		},
	})
}
