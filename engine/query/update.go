package query

import (
	"fmt"
	"strings"
)

// ============================================================================
// UPDATE BUILDER
// ============================================================================

// UpdateBuilder constrói queries UPDATE
type UpdateBuilder struct {
	Table        string
	Sets         []string
	SetValues    []interface{}
	WhereClauses []string
	WhereValues  []interface{}
}

// NewUpdate cria um novo UpdateBuilder
func NewUpdate(table string) *UpdateBuilder {
	return &UpdateBuilder{
		Table:        table,
		Sets:         []string{},
		SetValues:    []interface{}{},
		WhereClauses: []string{},
		WhereValues:  []interface{}{},
	}
}

// Set adiciona uma coluna e valor para atualizar
func (u *UpdateBuilder) Set(col string, val interface{}) *UpdateBuilder {
	u.Sets = append(u.Sets, fmt.Sprintf("%s = ?", col))
	u.SetValues = append(u.SetValues, val)
	return u
}

// Where adiciona condição WHERE
func (u *UpdateBuilder) Where(condition string, args ...interface{}) *UpdateBuilder {
	u.WhereClauses = append(u.WhereClauses, condition)
	u.WhereValues = append(u.WhereValues, args...)
	return u
}

// WhereRaw adiciona filtro customizado sem parâmetros
func (u *UpdateBuilder) WhereRaw(raw string) *UpdateBuilder {
	u.WhereClauses = append(u.WhereClauses, "("+raw+")")
	return u
}

// Build gera a query SQL final
func (u *UpdateBuilder) Build() (string, []interface{}) {
	setPart := strings.Join(u.Sets, ", ")
	wherePart := strings.Join(u.WhereClauses, " AND ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", u.Table, setPart, wherePart)

	// Combina valores do SET com valores do WHERE
	allValues := append(u.SetValues, u.WhereValues...)

	return query, allValues
}