package utils

import (
	"time"

	"github.com/Masterminds/squirrel"
)

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

func AddMonths(period string, monthsToAdd int) (string, error) {
	t, err := time.Parse("200601", period)

	if err != nil {
		return "", err
	}

	t = t.AddDate(0, monthsToAdd, 0)

	return t.Format("200601"), nil
}
