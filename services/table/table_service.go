package table

import (
	"database/sql"
	"fmt"
	"strings"
	"meu-provedor/config"
	"meu-provedor/models"
)

// ============================================================================
// TABLE OPERATIONS
// ============================================================================

// Create cria uma nova tabela para o projeto (usando project_id)
func Create(projectID int64, req models.CreateTableRequest) (string, error) {
	// Busca o code do projeto pelo ID
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return "", fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTableName := fmt.Sprintf("%s_%s", projectCode, req.TableName)

	columns := []string{
		"id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY",
		`id_instancia BIGINT UNSIGNED NOT NULL,
		FOREIGN KEY (id_instancia)
		REFERENCES instancias_projetion(id)
		ON DELETE CASCADE`,
	}

	for _, col := range req.Columns {
		def := col.Name + " " + col.Type
		if !col.Nullable {
			def += " NOT NULL"
		}
		if col.Unique {
			def += " UNIQUE"
		}
		columns = append(columns, def)
	}

	for _, idx := range req.Indexes {
		idxDef := ""
		if idx.Type == "UNIQUE" {
			idxDef = fmt.Sprintf("UNIQUE KEY %s (%s)", idx.Name, strings.Join(idx.Columns, ","))
		} else {
			idxDef = fmt.Sprintf("INDEX %s (%s)", idx.Name, strings.Join(idx.Columns, ","))
		}
		columns = append(columns, idxDef)
	}

	createSQL := fmt.Sprintf("CREATE TABLE %s (%s)", fullTableName, strings.Join(columns, ","))

	_, err = config.MasterDB.Exec(createSQL)
	return fullTableName, err
}

// List retorna todas as tabelas de um projeto (usando project_id)
func List(projectID int64) ([]string, error) {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return nil, fmt.Errorf("projeto não encontrado: %w", err)
	}

	rows, err := config.MasterDB.Query(`
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = DATABASE()
		AND table_name LIKE ?`, projectCode+"_%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var fullName string
		rows.Scan(&fullName)
		displayName := strings.TrimPrefix(fullName, projectCode+"_")
		tables = append(tables, displayName)
	}
	return tables, nil
}

// Delete remove uma tabela (usando project_id)
func Delete(projectID int64, table string) error {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, table)
	_, err = config.MasterDB.Exec("DROP TABLE " + fullTable)
	return err
}

// GetDetails retorna detalhes completos de uma tabela (usando project_id)
func GetDetails(projectID int64, tableName string) (*models.TableDetail, error) {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return nil, fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, tableName)

	columns, err := getColumns(fullTable)
	if err != nil {
		return nil, err
	}

	indexes, err := getIndexes(fullTable)
	if err != nil {
		return nil, err
	}

	return &models.TableDetail{
		Name:    tableName,
		Columns: columns,
		Indexes: indexes,
	}, nil
}

// ============================================================================
// COLUMN OPERATIONS
// ============================================================================

type ColumnRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Unique   bool   `json:"unique"`
}

// AddColumn adiciona uma nova coluna à tabela (usando project_id)
func AddColumn(projectID int64, tableName string, col ColumnRequest) error {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, tableName)
	def := col.Name + " " + col.Type
	if !col.Nullable {
		def += " NOT NULL"
	}
	if col.Unique {
		def += " UNIQUE"
	}
	query := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", fullTable, def)
	_, err = config.MasterDB.Exec(query)
	return err
}

// ModifyColumn modifica uma coluna existente (usando project_id)
func ModifyColumn(projectID int64, tableName string, col ColumnRequest) error {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, tableName)
	def := col.Name + " " + col.Type
	if !col.Nullable {
		def += " NOT NULL"
	}
	query := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s", fullTable, def)
	_, err = config.MasterDB.Exec(query)
	return err
}

// DropColumn remove uma coluna (usando project_id)
func DropColumn(projectID int64, tableName, columnName string) error {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, tableName)
	query := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", fullTable, columnName)
	_, err = config.MasterDB.Exec(query)
	return err
}

// ============================================================================
// INDEX OPERATIONS
// ============================================================================

type IndexRequest struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Type    string   `json:"type"`
}

// AddIndex adiciona um novo índice (usando project_id)
func AddIndex(projectID int64, tableName string, idx IndexRequest) error {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, tableName)
	var query string
	if idx.Type == "UNIQUE" {
		query = fmt.Sprintf("ALTER TABLE %s ADD UNIQUE INDEX %s (%s)",
			fullTable, idx.Name, strings.Join(idx.Columns, ","))
	} else {
		query = fmt.Sprintf("ALTER TABLE %s ADD INDEX %s (%s)",
			fullTable, idx.Name, strings.Join(idx.Columns, ","))
	}
	_, err = config.MasterDB.Exec(query)
	return err
}

// DropIndex remove um índice (usando project_id)
func DropIndex(projectID int64, tableName, indexName string) error {
	projectCode, err := config.GetProjectCodeByID(int(projectID))
	if err != nil {
		return fmt.Errorf("projeto não encontrado: %w", err)
	}

	fullTable := fmt.Sprintf("%s_%s", projectCode, tableName)
	query := fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", fullTable, indexName)
	_, err = config.MasterDB.Exec(query)
	return err
}

// ============================================================================
// INTERNAL HELPERS
// ============================================================================

func getColumns(fullTable string) ([]models.ColumnDetail, error) {
	rows, err := config.MasterDB.Query(`
		SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_DEFAULT, COLUMN_KEY, EXTRA
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`, fullTable,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []models.ColumnDetail
	for rows.Next() {
		var col models.ColumnDetail
		var isNullable, colDefault, colKey, extra sql.NullString

		err := rows.Scan(&col.Name, &col.Type, &isNullable, &colDefault, &colKey, &extra)
		if err != nil {
			continue
		}

		col.Nullable = (isNullable.String == "YES")
		if colDefault.Valid {
			col.Default = colDefault.String
		}
		col.Key = colKey.String
		col.Extra = extra.String

		columns = append(columns, col)
	}
	return columns, nil
}

func getIndexes(fullTable string) ([]models.IndexDetail, error) {
	rows, err := config.MasterDB.Query(`
		SELECT INDEX_NAME, COLUMN_NAME, NON_UNIQUE
		FROM information_schema.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`, fullTable,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexMap := make(map[string]*models.IndexDetail)
	for rows.Next() {
		var name, column string
		var nonUnique int

		err := rows.Scan(&name, &column, &nonUnique)
		if err != nil {
			continue
		}

		if _, exists := indexMap[name]; !exists {
			idxType := "INDEX"
			if nonUnique == 0 {
				idxType = "UNIQUE"
			}
			indexMap[name] = &models.IndexDetail{
				Name:    name,
				Columns: []string{},
				Type:    idxType,
			}
		}
		indexMap[name].Columns = append(indexMap[name].Columns, column)
	}

	var indexes []models.IndexDetail
	for _, idx := range indexMap {
		indexes = append(indexes, *idx)
	}
	return indexes, nil
}
