package query

import (
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// SOFT DELETE BUILDER (UPDATE deleted_at)
// ============================================================================

// SoftDeleteBuilder constrói queries de soft delete (UPDATE deleted_at)
type SoftDeleteBuilder struct {
	Table    string
	Where    []string
	RawWhere []string
	Values   []interface{}
}

// NewSoftDelete cria um novo SoftDeleteBuilder
func NewSoftDelete(table string) *SoftDeleteBuilder {
	return &SoftDeleteBuilder{
		Table:    table,
		Where:    []string{},
		RawWhere: []string{},
		Values:   []interface{}{},
	}
}

// AddWhere adiciona condição WHERE parametrizada
func (d *SoftDeleteBuilder) AddWhere(condition string, args ...interface{}) *SoftDeleteBuilder {
	d.Where = append(d.Where, condition)
	d.Values = append(d.Values, args...)
	return d
}

// AddRawWhere adiciona condição WHERE sem parâmetros
func (d *SoftDeleteBuilder) AddRawWhere(condition string) *SoftDeleteBuilder {
	d.RawWhere = append(d.RawWhere, condition)
	return d
}

// Build gera a query SQL final
func (d *SoftDeleteBuilder) Build(deletedAt time.Time) (string, []interface{}) {
	where := []string{}
	where = append(where, d.Where...)
	where = append(where, d.RawWhere...)

	query := fmt.Sprintf("UPDATE %s SET deleted_at = ?", d.Table)

	// deleted_at é o primeiro valor
	args := []interface{}{deletedAt}
	args = append(args, d.Values...)

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	return query, args
}