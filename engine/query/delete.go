package query

import (
	"fmt"
	"strings"
)

// ============================================================================
// DELETE BUILDER
// ============================================================================

// DeleteBuilder constrói queries DELETE
type DeleteBuilder struct {
	Table        string
	WhereClauses []string
	WhereValues  []interface{}
}

// NewDelete cria um novo DeleteBuilder
func NewDelete(table string) *DeleteBuilder {
	return &DeleteBuilder{
		Table:        table,
		WhereClauses: []string{},
		WhereValues:  []interface{}{},
	}
}

// Where adiciona condição WHERE
func (d *DeleteBuilder) Where(condition string, args ...interface{}) *DeleteBuilder {
	d.WhereClauses = append(d.WhereClauses, condition)
	d.WhereValues = append(d.WhereValues, args...)
	return d
}

// WhereRaw adiciona filtro customizado sem parâmetros
func (d *DeleteBuilder) WhereRaw(raw string) *DeleteBuilder {
	d.WhereClauses = append(d.WhereClauses, "("+raw+")")
	return d
}

// Build gera a query SQL final
func (d *DeleteBuilder) Build() (string, []interface{}) {
	wherePart := strings.Join(d.WhereClauses, " AND ")
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", d.Table, wherePart)
	return query, d.WhereValues
}