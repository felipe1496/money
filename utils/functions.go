package utils

import (
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
)

func ErrorResponse(status int, msg string) gin.H {
	return gin.H{
		"status": status,
		"error": gin.H{
			"message": msg,
			"type":    StatusMessages[status],
		},
	}
}

func ApplyFilterToSquirrel(query squirrel.SelectBuilder, filter *FilterBuilder) (squirrel.SelectBuilder, error) {
	if filter == nil {
		return query, nil
	}

	if filter.HasError() {
		return query, filter.GetError()
	}

	sql, args, err := filter.Build()
	if err != nil {
		return query, err
	}
	if sql != "" {
		query = query.Where(squirrel.Expr(sql, args...))
	}

	orderBy := filter.GetOrderBy()
	for _, order := range orderBy {
		query = query.OrderBy(order)
	}

	if limit := filter.GetLimit(); limit != nil {
		query = query.Limit(*limit)
	}

	if offset := filter.GetOffset(); offset != nil {
		query = query.Offset(*offset)
	}

	return query, nil
}
