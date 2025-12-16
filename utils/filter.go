package utils

import (
	"fmt"
	"strings"
)

type Condition struct {
	Field    string
	Operator string
	Value    any
}

type OrderByClause struct {
	Field     string
	Direction string
}

type FilterBuilder struct {
	conditions []any
	args       []any
	orderBy    []OrderByClause
	limitVal   *int
	offsetVal  *int
	err        error
}

var operatorMap = map[string]string{
	"eq":    "=",
	"ne":    "!=",
	"gt":    ">",
	"gte":   ">=",
	"lt":    "<",
	"lte":   "<=",
	"like":  "LIKE",
	"nlike": "NOT LIKE",
	"in":    "IN",
	"nin":   "NOT IN",
	"is":    "IS",
	"isn":   "IS NOT",
}

func CreateFilter() *FilterBuilder {
	return &FilterBuilder{
		conditions: make([]any, 0),
		args:       make([]any, 0),
		orderBy:    make([]OrderByClause, 0),
	}
}

func (fb *FilterBuilder) And(field, operator string, value any) *FilterBuilder {
	if fb.err != nil {
		return fb
	}

	if !isValidOperator(operator) {
		fb.err = fmt.Errorf("invalid operator: %s", operator)
		return fb
	}

	fb.conditions = append(fb.conditions, Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
	})

	return fb
}

func (fb *FilterBuilder) Or(conditions ...Condition) *FilterBuilder {
	if fb.err != nil {
		return fb
	}

	for _, cond := range conditions {
		if !isValidOperator(cond.Operator) {
			fb.err = fmt.Errorf("invalid operator: %s", cond.Operator)
			return fb
		}
	}

	if len(conditions) > 0 {
		fb.conditions = append(fb.conditions, conditions)
	}

	return fb
}

func (fb *FilterBuilder) OrderBy(field string, direction string) *FilterBuilder {
	if fb.err != nil {
		return fb
	}

	direction = strings.ToUpper(direction)

	if direction != "ASC" && direction != "DESC" {
		fb.err = fmt.Errorf("direção inválida para ORDER BY: %s (use 'asc' ou 'desc')", direction)
		return fb
	}

	if strings.TrimSpace(field) == "" {
		fb.err = fmt.Errorf("campo vazio para ORDER BY")
		return fb
	}

	fb.orderBy = append(fb.orderBy, OrderByClause{
		Field:     field,
		Direction: direction,
	})

	return fb
}

func (fb *FilterBuilder) Limit(limit int) *FilterBuilder {
	if fb.err != nil {
		return fb
	}

	if limit < 0 {
		fb.err = fmt.Errorf("limit não pode ser negativo: %d", limit)
		return fb
	}

	fb.limitVal = &limit
	return fb
}

func (fb *FilterBuilder) Offset(offset int) *FilterBuilder {
	if fb.err != nil {
		return fb
	}

	if offset < 0 {
		fb.err = fmt.Errorf("offset não pode ser negativo: %d", offset)
		return fb
	}

	fb.offsetVal = &offset
	return fb
}

func (fb *FilterBuilder) Build() (string, []any, error) {
	if fb.err != nil {
		return "", nil, fb.err
	}

	if len(fb.conditions) == 0 {
		return "", []any{}, nil
	}

	var whereParts []string
	argIndex := 1

	for _, cond := range fb.conditions {
		switch c := cond.(type) {
		case Condition:

			sqlOp := operatorMap[c.Operator]
			whereParts = append(whereParts, fmt.Sprintf("%s %s $%d", c.Field, sqlOp, argIndex))
			fb.args = append(fb.args, c.Value)
			argIndex++

		case []Condition:

			var orParts []string
			for _, orCond := range c {
				sqlOp := operatorMap[orCond.Operator]
				orParts = append(orParts, fmt.Sprintf("%s %s $%d", orCond.Field, sqlOp, argIndex))
				fb.args = append(fb.args, orCond.Value)
				argIndex++
			}
			whereParts = append(whereParts, fmt.Sprintf("(%s)", strings.Join(orParts, " OR ")))
		}
	}

	whereClause := strings.Join(whereParts, " AND ")
	return whereClause, fb.args, nil
}

func (fb *FilterBuilder) GetOrderBy() []string {
	if len(fb.orderBy) == 0 {
		return nil
	}

	result := make([]string, len(fb.orderBy))
	for i, order := range fb.orderBy {
		result[i] = fmt.Sprintf("%s %s", order.Field, order.Direction)
	}
	return result
}

func (fb *FilterBuilder) GetLimit() *uint64 {
	if fb.limitVal == nil {
		return nil
	}
	val := uint64(*fb.limitVal)
	return &val
}

func (fb *FilterBuilder) GetOffset() *uint64 {
	if fb.offsetVal == nil {
		return nil
	}
	val := uint64(*fb.offsetVal)
	return &val
}

func (fb *FilterBuilder) HasError() bool {
	return fb.err != nil
}

func (fb *FilterBuilder) GetError() error {
	return fb.err
}

func isValidOperator(op string) bool {
	_, exists := operatorMap[op]
	return exists
}

func GetValidOperators() []string {
	ops := make([]string, 0, len(operatorMap))
	for k := range operatorMap {
		ops = append(ops, k)
	}
	return ops
}
