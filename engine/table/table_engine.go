package table

import (
	"fmt"
	"strings"

	"meu-provedor/config"
	"meu-provedor/models"
)

func CreateTable(projectCode string, req models.CreateTableRequest) (string, error) {
	fullTableName := fmt.Sprintf("%s_%s", projectCode, req.TableName)

	columns := []string{
		"id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY",
		`id_instancia BIGINT UNSIGNED NOT NULL,
			FOREIGN KEY (id_instancia)
			REFERENCES instancias_projetion(id)
			ON DELETE CASCADE`,
	}

	// colunas personalizadas
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

	// índices
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

	_, err := config.MasterDB.Exec(createSQL)
	return fullTableName, err
}

func ListTables(projectCode string) ([]string, error) {
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

func DropTable(projectCode, table string) error {
	fullTable := fmt.Sprintf("%s_%s", projectCode, table)
	_, err := config.MasterDB.Exec("DROP TABLE " + fullTable)
	return err
}


