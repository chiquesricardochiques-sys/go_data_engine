package query

import (
	"fmt"
	"strings"
)

// ============================================================================
// JOIN SELECT BUILDER (Para queries complexas com múltiplos JOINs)
// ============================================================================

// JoinConfig configuração de um JOIN
type JoinConfig struct {
	Type    string
	Table   string
	Alias   string
	On      string
	Columns []string
}

// JoinSelectBuilder constrói queries SELECT com múltiplos JOINs
type JoinSelectBuilder struct {
	BaseTable string
	BaseAlias string
	Columns   []string
	Joins     []JoinConfig
	Where     []string
	RawWhere  []string
	GroupBy   string
	Having    string
	OrderBy   string
	Limit     int
	Offset    int
	Values    []interface{}
}

// NewJoinSelect cria um novo JoinSelectBuilder
func NewJoinSelect(table, alias string) *JoinSelectBuilder {
	if alias == "" {
		alias = table
	}
	return &JoinSelectBuilder{
		BaseTable: table,
		BaseAlias: alias,
		Columns:   []string{},
		Joins:     []JoinConfig{},
		Where:     []string{},
		RawWhere:  []string{},
		Values:    []interface{}{},
	}
}

// AddColumns adiciona colunas ao SELECT
func (b *JoinSelectBuilder) AddColumns(cols ...string) *JoinSelectBuilder {
	b.Columns = append(b.Columns, cols...)
	return b
}

// AddJoin adiciona uma configuração de JOIN
func (b *JoinSelectBuilder) AddJoin(j JoinConfig) *JoinSelectBuilder {
	j.Type = NormalizeJoinType(j.Type)
	if j.Alias == "" {
		j.Alias = j.Table
	}
	b.Joins = append(b.Joins, j)
	return b
}

// AddWhere adiciona condição WHERE parametrizada
func (b *JoinSelectBuilder) AddWhere(cond string, args ...interface{}) *JoinSelectBuilder {
	b.Where = append(b.Where, cond)
	b.Values = append(b.Values, args...)
	return b
}

// AddRawWhere adiciona condição WHERE sem parâmetros
func (b *JoinSelectBuilder) AddRawWhere(cond string) *JoinSelectBuilder {
	b.RawWhere = append(b.RawWhere, cond)
	return b
}

// SetGroupBy define GROUP BY
func (b *JoinSelectBuilder) SetGroupBy(group string) *JoinSelectBuilder {
	b.GroupBy = group
	return b
}

// SetHaving define HAVING
func (b *JoinSelectBuilder) SetHaving(having string) *JoinSelectBuilder {
	b.Having = having
	return b
}

// SetOrderBy define ORDER BY
func (b *JoinSelectBuilder) SetOrderBy(order string) *JoinSelectBuilder {
	b.OrderBy = order
	return b
}

// SetLimitOffset define LIMIT e OFFSET
func (b *JoinSelectBuilder) SetLimitOffset(limit, offset int) *JoinSelectBuilder {
	b.Limit = limit
	b.Offset = offset
	return b
}

// Build gera a query SQL final
func (b *JoinSelectBuilder) Build() (string, []interface{}) {
	// Se não há colunas especificadas, usa *
	if len(b.Columns) == 0 {
		b.Columns = append(b.Columns, "*")
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s AS %s",
		strings.Join(b.Columns, ", "),
		b.BaseTable,
		b.BaseAlias,
	)

	// Adiciona JOINs
	for _, j := range b.Joins {
		query += fmt.Sprintf(
			" %s JOIN %s AS %s ON %s",
			j.Type,
			j.Table,
			j.Alias,
			j.On,
		)
	}

	// Combina WHERE parametrizado e raw
	where := []string{}
	where = append(where, b.Where...)
	where = append(where, b.RawWhere...)

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	if b.GroupBy != "" {
		query += " GROUP BY " + b.GroupBy
	}

	if b.Having != "" {
		query += " HAVING " + b.Having
	}

	if b.OrderBy != "" {
		query += " ORDER BY " + b.OrderBy
	}

	if b.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", b.Limit)
		if b.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", b.Offset)
		}
	}

	return query, b.Values
}