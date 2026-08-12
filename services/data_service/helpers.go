package services

import (
	"database/sql"
	"fmt"
	"meu-provedor/config"
	"meu-provedor/models"
)

// ============================================================================
// HELPER FUNCTIONS - Funções compartilhadas entre services
// ============================================================================

// GetProjectCodeByID busca o código do projeto pelo ID
func GetProjectCodeByID(projectID int64) (string, error) {
	var code string
	query := "SELECT code FROM projects WHERE id = ? LIMIT 1"
	
	err := config.MasterDB.QueryRow(query, projectID).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", models.ErrProjectNotFound
		}
		return "", fmt.Errorf("erro ao buscar projeto: %w", err)
	}
	
	if code == "" {
		return "", models.ErrProjectNotFound
	}
	
	return code, nil
}

// BuildTableName constrói o nome físico da tabela com prefixo do projeto
func BuildTableName(projectCode, table string) (string, error) {
	if table == "" {
		return "", models.ErrTableRequired
	}
	return fmt.Sprintf("%s_%s", projectCode, table), nil
}

// RowsToMap converte sql.Rows para []map[string]interface{}
func RowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	
	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		
		for i := range values {
			ptrs[i] = &values[i]
		}
		
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		
		row := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		
		result = append(result, row)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return result, nil
}

// EnsureSoftDeleteColumn garante que a coluna deleted_at existe na tabela
func EnsureSoftDeleteColumn(db *sql.DB, table string) error {
	var exists int
	queryCheck := `
		SELECT COUNT(*)
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		  AND TABLE_NAME = ?
		  AND COLUMN_NAME = 'deleted_at'
	`
	
	if err := db.QueryRow(queryCheck, table).Scan(&exists); err != nil {
		return fmt.Errorf("erro ao verificar coluna deleted_at: %w", err)
	}
	
	if exists == 0 {
		alter := fmt.Sprintf("ALTER TABLE %s ADD COLUMN deleted_at DATETIME NULL", table)
		if _, err := db.Exec(alter); err != nil {
			return fmt.Errorf("erro ao criar coluna deleted_at: %w", err)
		}
	}
	
	return nil
}

