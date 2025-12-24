package middlewares

import (
	"rango-backend/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func QueryOptsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, _ := ctx.GetQuery("page")
		perPage, _ := ctx.GetQuery("per_page")
		orderBy, _ := ctx.GetQuery("order_by")
		order, _ := ctx.GetQuery("order")
		pageNum, err := strconv.Atoi(page)

		if err != nil {
			pageNum = 1
		}

		perPageNum, err := strconv.Atoi(perPage)

		if err != nil {
			perPageNum = 10
		}

		queryOpts := utils.QueryOpts()
		queryOpts.Offset((pageNum - 1) * perPageNum)
		queryOpts.Limit(perPageNum + 1)

		if orderBy != "" && order != "" {
			queryOpts.OrderBy(orderBy, order)
		}

		ctx.Set("page", pageNum)
		ctx.Set("per_page", perPageNum)
		ctx.Set("query_opts", queryOpts)
		ctx.Next()
	}
}
