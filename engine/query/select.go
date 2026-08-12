package query

import (
	"fmt"
	"strings"
)

// ============================================================================
// SELECT BUILDER
// ============================================================================

// SelectBuilder constrói queries SELECT
type SelectBuilder struct {
	Table   string
	Alias   string
	Columns []string
	Joins   []string
	Where   []string
	GroupBy string
	Having  string
	OrderBy string
	Limit   int
	Offset  int
	Values  []interface{}
}

// NewSelect cria um novo SelectBuilder
func NewSelect(table, alias string) *SelectBuilder {
	if alias == "" {
		alias = table
	}
	return &SelectBuilder{
		Table:   table,
		Alias:   alias,
		Columns: []string{"*"},
		Joins:   []string{},
		Where:   []string{},
		Values:  []interface{}{},
	}
}

// SetColumns define as colunas a serem selecionadas
func (s *SelectBuilder) SetColumns(cols []string) *SelectBuilder {
	if len(cols) > 0 {
		s.Columns = cols
	}
	return s
}

// AddJoin adiciona uma cláusula JOIN
func (s *SelectBuilder) AddJoin(joinType, table, alias, on string) *SelectBuilder {
	joinType = NormalizeJoinType(joinType)
	if alias == "" {
		alias = table
	}
	joinClause := fmt.Sprintf("%s JOIN %s AS %s ON %s", joinType, table, alias, on)
	s.Joins = append(s.Joins, joinClause)
	return s
}

// AddWhere adiciona condição WHERE
func (s *SelectBuilder) AddWhere(condition string, args ...interface{}) *SelectBuilder {
	s.Where = append(s.Where, condition)
	s.Values = append(s.Values, args...)
	return s
}

// SetGroupBy define cláusula GROUP BY
func (s *SelectBuilder) SetGroupBy(group string) *SelectBuilder {
	s.GroupBy = group
	return s
}

// SetHaving define cláusula HAVING
func (s *SelectBuilder) SetHaving(having string) *SelectBuilder {
	s.Having = having
	return s
}

// SetOrderBy define cláusula ORDER BY
func (s *SelectBuilder) SetOrderBy(order string) *SelectBuilder {
	s.OrderBy = order
	return s
}

// SetLimitOffset define LIMIT e OFFSET
func (s *SelectBuilder) SetLimitOffset(limit, offset int) *SelectBuilder {
	s.Limit = limit
	s.Offset = offset
	return s
}

// Build gera a query SQL final
func (s *SelectBuilder) Build() string {
	query := fmt.Sprintf("SELECT %s FROM %s AS %s",
		strings.Join(s.Columns, ", "),
		s.Table,
		s.Alias,
	)

	if len(s.Joins) > 0 {
		query += " " + strings.Join(s.Joins, " ")
	}

	if len(s.Where) > 0 {
		query += " WHERE " + strings.Join(s.Where, " AND ")
	}

	if s.GroupBy != "" {
		query += " GROUP BY " + s.GroupBy
		if s.Having != "" {
			query += " HAVING " + s.Having
		}
	}

	if s.OrderBy != "" {
		query += " ORDER BY " + s.OrderBy
	}

	if s.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", s.Limit)
		if s.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", s.Offset)
		}
	}

	return query
}

// GetValues retorna os valores dos parâmetros
func (s *SelectBuilder) GetValues() []interface{} {
	return s.Values
}