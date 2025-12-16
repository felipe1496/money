package tests

import (
	"rango-backend/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFilter(t *testing.T) {
	filter := utils.CreateFilter()
	assert.NotNil(t, filter)
}

func TestSimpleAndCondition(t *testing.T) {
	filter := utils.CreateFilter().And("id", "eq", 10)
	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id = $1", sql)
	assert.Equal(t, []interface{}{10}, args)
}

func TestMultipleAndConditions(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 10).
		And("name", "like", "%test%").
		And("age", "gte", 18)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id = $1 AND name LIKE $2 AND age >= $3", sql)
	assert.Equal(t, []interface{}{10, "%test%", 18}, args)
}

func TestOrConditions(t *testing.T) {
	filter := utils.CreateFilter().Or(
		utils.Condition{Field: "status", Operator: "eq", Value: "active"},
		utils.Condition{Field: "status", Operator: "eq", Value: "pending"},
	)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "(status = $1 OR status = $2)", sql)
	assert.Equal(t, []interface{}{"active", "pending"}, args)
}

func TestMixedAndOrConditions(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "gt", 5).
		Or(
			utils.Condition{Field: "status", Operator: "ne", Value: "inactive"},
			utils.Condition{Field: "status", Operator: "ne", Value: "deleted"},
		).
		And("created_at", "lt", "2024-01-01")

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id > $1 AND (status != $2 OR status != $3) AND created_at < $4", sql)
	assert.Equal(t, []interface{}{5, "inactive", "deleted", "2024-01-01"}, args)
}

func TestComplexExample(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 10).
		And("name", "nlike", "%test%").
		Or(
			utils.Condition{Field: "name", Operator: "ne", Value: "felipe"},
			utils.Condition{Field: "name", Operator: "ne", Value: "roberto"},
		)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id = $1 AND name NOT LIKE $2 AND (name != $3 OR name != $4)", sql)
	assert.Equal(t, []interface{}{10, "%test%", "felipe", "roberto"}, args)
}

func TestInvalidOperator(t *testing.T) {
	filter := utils.CreateFilter().And("id", "invalid", 10)
	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid operator: invalid")
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestInvalidOperatorInOr(t *testing.T) {
	filter := utils.CreateFilter().Or(
		utils.Condition{Field: "status", Operator: "invalid_op", Value: "active"},
	)

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid operator: invalid_op")
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestEmptyFilter(t *testing.T) {
	filter := utils.CreateFilter()
	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Empty(t, sql)
	assert.Empty(t, args)
}

func TestAllOperators(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected string
	}{
		{"equal", "eq", "="},
		{"not equal", "ne", "!="},
		{"greater than", "gt", ">"},
		{"greater than or equal", "gte", ">="},
		{"less than", "lt", "<"},
		{"less than or equal", "lte", "<="},
		{"like", "like", "LIKE"},
		{"not like", "nlike", "NOT LIKE"},
		{"in", "in", "IN"},
		{"not in", "nin", "NOT IN"},
		{"is", "is", "IS"},
		{"is not", "isn", "IS NOT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := utils.CreateFilter().And("field", tt.operator, "value")
			sql, args, err := filter.Build()
			require.NoError(t, err)
			assert.Equal(t, "field "+tt.expected+" $1", sql)
			assert.Equal(t, []interface{}{"value"}, args)
		})
	}
}

func TestChainedErrorPropagation(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 10).
		And("name", "invalid_op", "test").
		And("age", "gt", 18)

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid operator: invalid_op")
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestMultipleOrGroups(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "gt", 0).
		Or(
			utils.Condition{Field: "type", Operator: "eq", Value: "admin"},
			utils.Condition{Field: "type", Operator: "eq", Value: "moderator"},
		).
		Or(
			utils.Condition{Field: "active", Operator: "eq", Value: true},
			utils.Condition{Field: "verified", Operator: "eq", Value: true},
		)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id > $1 AND (type = $2 OR type = $3) AND (active = $4 OR verified = $5)", sql)
	assert.Equal(t, []interface{}{0, "admin", "moderator", true, true}, args)
}

func TestSingleOrCondition(t *testing.T) {
	filter := utils.CreateFilter().Or(
		utils.Condition{Field: "status", Operator: "eq", Value: "active"},
	)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "(status = $1)", sql)
	assert.Equal(t, []interface{}{"active"}, args)
}

func TestEmptyOrConditions(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		Or()

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id = $1", sql)
	assert.Equal(t, []interface{}{1}, args)
}

func TestDifferentValueTypes(t *testing.T) {
	filter := utils.CreateFilter().
		And("int_field", "eq", 42).
		And("string_field", "like", "test").
		And("bool_field", "eq", true).
		And("float_field", "gte", 3.14)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "int_field = $1 AND string_field LIKE $2 AND bool_field = $3 AND float_field >= $4", sql)
	assert.Equal(t, []interface{}{42, "test", true, 3.14}, args)
}

func TestGetValidOperators(t *testing.T) {
	operators := utils.GetValidOperators()
	assert.NotEmpty(t, operators)
	assert.Contains(t, operators, "eq")
	assert.Contains(t, operators, "ne")
	assert.Contains(t, operators, "like")
	assert.Len(t, operators, 12)
}

// ===== TESTES PARA ORDER BY =====

func TestGetOrderBySimple(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "gt", 5).
		OrderBy("name", "asc")

	// Build retorna apenas WHERE
	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id > $1", sql)
	assert.Equal(t, []interface{}{5}, args)

	// GetOrderBy retorna ORDER BY separadamente
	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"name ASC"}, orderBy)
}

func TestGetOrderByDescending(t *testing.T) {
	filter := utils.CreateFilter().
		And("status", "eq", "active").
		OrderBy("created_at", "desc")

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "status = $1", sql)
	assert.Equal(t, []interface{}{"active"}, args)

	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"created_at DESC"}, orderBy)
}

func TestGetOrderByMultiple(t *testing.T) {
	filter := utils.CreateFilter().
		And("status", "eq", "active").
		OrderBy("priority", "desc").
		OrderBy("name", "asc").
		OrderBy("created_at", "desc")

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "status = $1", sql)
	assert.Equal(t, []interface{}{"active"}, args)

	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"priority DESC", "name ASC", "created_at DESC"}, orderBy)
}

func TestGetOrderByWithoutConditions(t *testing.T) {
	filter := utils.CreateFilter().
		OrderBy("id", "asc")

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Empty(t, sql)
	assert.Empty(t, args)

	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"id ASC"}, orderBy)
}

func TestGetOrderByEmpty(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1)

	sql, _, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id = $1", sql)

	orderBy := filter.GetOrderBy()
	assert.Nil(t, orderBy)
}

func TestOrderByInvalidDirection(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		OrderBy("name", "invalid")

	assert.True(t, filter.HasError())
	assert.Error(t, filter.GetError())
	assert.Contains(t, filter.GetError().Error(), "direção inválida para ORDER BY")

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "direção inválida para ORDER BY")
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestOrderByEmptyField(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		OrderBy("", "asc")

	assert.True(t, filter.HasError())
	assert.Error(t, filter.GetError())
	assert.Contains(t, filter.GetError().Error(), "campo vazio para ORDER BY")

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestOrderByCaseInsensitive(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		OrderBy("name", "AsC")

	sql, _, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id = $1", sql)

	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"name ASC"}, orderBy)
}

// ===== TESTES PARA LIMIT =====

func TestGetLimitSimple(t *testing.T) {
	filter := utils.CreateFilter().
		And("status", "eq", "active").
		Limit(10)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "status = $1", sql)
	assert.Equal(t, []interface{}{"active"}, args)

	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(10), *limit)
}

func TestGetLimitWithoutConditions(t *testing.T) {
	filter := utils.CreateFilter().
		Limit(5)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Empty(t, sql)
	assert.Empty(t, args)

	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(5), *limit)
}

func TestGetLimitZero(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "gt", 0).
		Limit(0)

	sql, _, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id > $1", sql)

	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(0), *limit)
}

func TestGetLimitNotSet(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1)

	limit := filter.GetLimit()
	assert.Nil(t, limit)
}

func TestLimitNegative(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		Limit(-1)

	assert.True(t, filter.HasError())
	assert.Error(t, filter.GetError())
	assert.Contains(t, filter.GetError().Error(), "limit não pode ser negativo")

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limit não pode ser negativo")
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

// ===== TESTES PARA OFFSET =====

func TestGetOffsetSimple(t *testing.T) {
	filter := utils.CreateFilter().
		And("status", "eq", "active").
		Offset(20)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "status = $1", sql)
	assert.Equal(t, []interface{}{"active"}, args)

	offset := filter.GetOffset()
	require.NotNil(t, offset)
	assert.Equal(t, uint64(20), *offset)
}

func TestGetOffsetWithoutConditions(t *testing.T) {
	filter := utils.CreateFilter().
		Offset(10)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Empty(t, sql)
	assert.Empty(t, args)

	offset := filter.GetOffset()
	require.NotNil(t, offset)
	assert.Equal(t, uint64(10), *offset)
}

func TestGetOffsetZero(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "gt", 0).
		Offset(0)

	sql, _, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "id > $1", sql)

	offset := filter.GetOffset()
	require.NotNil(t, offset)
	assert.Equal(t, uint64(0), *offset)
}

func TestGetOffsetNotSet(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1)

	offset := filter.GetOffset()
	assert.Nil(t, offset)
}

func TestOffsetNegative(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		Offset(-5)

	assert.True(t, filter.HasError())
	assert.Error(t, filter.GetError())
	assert.Contains(t, filter.GetError().Error(), "offset não pode ser negativo")

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "offset não pode ser negativo")
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

// ===== TESTES COMBINADOS =====

func TestGetLimitAndOffset(t *testing.T) {
	filter := utils.CreateFilter().
		And("status", "eq", "active").
		Limit(10).
		Offset(20)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "status = $1", sql)
	assert.Equal(t, []interface{}{"active"}, args)

	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(10), *limit)

	offset := filter.GetOffset()
	require.NotNil(t, offset)
	assert.Equal(t, uint64(20), *offset)
}

func TestGetOrderByWithLimit(t *testing.T) {
	filter := utils.CreateFilter().
		And("category", "eq", "tech").
		OrderBy("created_at", "desc").
		Limit(5)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "category = $1", sql)
	assert.Equal(t, []interface{}{"tech"}, args)

	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"created_at DESC"}, orderBy)

	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(5), *limit)
}

func TestCompleteQueryWithAllFeatures(t *testing.T) {
	filter := utils.CreateFilter().
		And("status", "eq", "active").
		And("verified", "eq", true).
		Or(
			utils.Condition{Field: "role", Operator: "eq", Value: "admin"},
			utils.Condition{Field: "role", Operator: "eq", Value: "moderator"},
		).
		OrderBy("priority", "desc").
		OrderBy("name", "asc").
		Limit(25).
		Offset(50)

	// Build retorna apenas WHERE
	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Equal(t, "status = $1 AND verified = $2 AND (role = $3 OR role = $4)", sql)
	assert.Equal(t, []interface{}{"active", true, "admin", "moderator"}, args)

	// Verificar ORDER BY
	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"priority DESC", "name ASC"}, orderBy)

	// Verificar LIMIT
	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(25), *limit)

	// Verificar OFFSET
	offset := filter.GetOffset()
	require.NotNil(t, offset)
	assert.Equal(t, uint64(50), *offset)
}

func TestGetOrderByLimitOffsetWithoutWhere(t *testing.T) {
	filter := utils.CreateFilter().
		OrderBy("id", "asc").
		Limit(10).
		Offset(5)

	sql, args, err := filter.Build()
	require.NoError(t, err)
	assert.Empty(t, sql)
	assert.Empty(t, args)

	orderBy := filter.GetOrderBy()
	assert.Equal(t, []string{"id ASC"}, orderBy)

	limit := filter.GetLimit()
	require.NotNil(t, limit)
	assert.Equal(t, uint64(10), *limit)

	offset := filter.GetOffset()
	require.NotNil(t, offset)
	assert.Equal(t, uint64(5), *offset)
}

func TestErrorPropagationWithOrderBy(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "invalid_op", 1).
		OrderBy("name", "asc")

	assert.True(t, filter.HasError())

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestErrorPropagationWithLimit(t *testing.T) {
	filter := utils.CreateFilter().
		And("id", "eq", 1).
		OrderBy("name", "invalid_direction").
		Limit(10)

	assert.True(t, filter.HasError())

	sql, args, err := filter.Build()
	assert.Error(t, err)
	assert.Empty(t, sql)
	assert.Nil(t, args)
}

func TestHasErrorAndGetError(t *testing.T) {
	// Sem erro
	filter1 := utils.CreateFilter().And("id", "eq", 1)
	assert.False(t, filter1.HasError())
	assert.Nil(t, filter1.GetError())

	// Com erro
	filter2 := utils.CreateFilter().And("id", "invalid", 1)
	assert.True(t, filter2.HasError())
	assert.NotNil(t, filter2.GetError())
	assert.Contains(t, filter2.GetError().Error(), "invalid operator: invalid")
}

func TestMultipleErrors(t *testing.T) {
	// Primeiro erro para a execução
	filter := utils.CreateFilter().
		And("id", "invalid_op", 1).
		OrderBy("", "asc").
		Limit(-1)

	// Apenas o primeiro erro é capturado
	assert.True(t, filter.HasError())
	assert.Contains(t, filter.GetError().Error(), "invalid operator: invalid_op")
}
