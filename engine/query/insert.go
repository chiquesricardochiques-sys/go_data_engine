package query

import (
	"fmt"
	"strings"
)

type InsertBuilder struct {
	table      string
	columns    []string
	values     [][]interface{}
	validated  bool
}

func NewInsert(table string) *InsertBuilder {
	return &InsertBuilder{
		table:     table,
		columns:   []string{},
		values:    [][]interface{}{},
		validated: false,
	}
}

// SetColumns define as colunas do INSERT
func (b *InsertBuilder) SetColumns(cols []string) *InsertBuilder {
	b.columns = cols
	return b
}

// AddRow adiciona uma linha de valores
func (b *InsertBuilder) AddRow(vals []interface{}) error {
	// Validar que número de valores = número de colunas
	if len(vals) != len(b.columns) {
		return fmt.Errorf("número de valores (%d) diferente de colunas (%d)", 
			len(vals), len(b.columns))
	}
	
	b.values = append(b.values, vals)
	return nil
}

// Build gera SQL MySQL com placeholders ?
func (b *InsertBuilder) Build() (string, []interface{}, error) {
	// Validações
	if b.table == "" {
		return "", nil, fmt.Errorf("tabela não definida")
	}
	if len(b.columns) == 0 {
		return "", nil, fmt.Errorf("nenhuma coluna definida")
	}
	if len(b.values) == 0 {
		return "", nil, fmt.Errorf("nenhum valor fornecido")
	}
	
	// Montar placeholders
	var placeholderGroups []string
	var allValues []interface{}
	
	for _, row := range b.values {
		// Para cada row: (?, ?, ?)
		placeholders := make([]string, len(b.columns))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		
		placeholderGroups = append(placeholderGroups, 
			"("+strings.Join(placeholders, ",")+")")
		
		allValues = append(allValues, row...)
	}
	
	// Montar query final
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		b.table,
		strings.Join(b.columns, ","),
		strings.Join(placeholderGroups, ","),
	)
	
	return query, allValues, nil
}
