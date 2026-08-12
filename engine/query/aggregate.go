package query

import (
	"fmt"
	"strings"
)

// ============================================================================
// AGGREGATE BUILDER (COUNT, SUM, AVG, MIN, MAX, EXISTS)
// ============================================================================

// AggregateBuilder constrói queries de agregação
type AggregateBuilder struct {
	Table     string
	Alias     string
	Operation string
	Column    string
	Where     []string
	Values    []interface{}
}

// NewAggregate cria um novo AggregateBuilder
func NewAggregate(table, alias, operation, column string) *AggregateBuilder {
	if alias == "" {
		alias = table
	}
	return &AggregateBuilder{
		Table:     table,
		Alias:     alias,
		Operation: NormalizeOperation(operation),
		Column:    column,
		Where:     []string{},
		Values:    []interface{}{},
	}
}

// AddWhere adiciona condição WHERE
func (a *AggregateBuilder) AddWhere(condition string, args ...interface{}) *AggregateBuilder {
	a.Where = append(a.Where, condition)
	a.Values = append(a.Values, args...)
	return a
}

// Build gera a query SQL final
func (a *AggregateBuilder) Build() string {
	target := "*"
	if a.Column != "" {
		target = a.Column
	}

	var selectExpr string
	if a.Operation == "EXISTS" {
		selectExpr = "EXISTS(SELECT 1"
	} else {
		selectExpr = fmt.Sprintf("%s(%s)", a.Operation, target)
	}

	query := fmt.Sprintf("SELECT %s FROM %s AS %s", selectExpr, a.Table, a.Alias)

	if len(a.Where) > 0 {
		query += " WHERE " + strings.Join(a.Where, " AND ")
	}

	if a.Operation == "EXISTS" {
		query += ")"
	}

	return query
}

// GetValues retorna os valores dos parâmetros
func (a *AggregateBuilder) GetValues() []interface{} {
	return a.Values
}