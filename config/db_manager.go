package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

// ============================================================================
// DATABASE CONNECTION
// ============================================================================

// MasterDB é a conexão global com o banco de dados
var MasterDB *sql.DB

// ConnectMaster estabelece conexão com o banco MySQL local
func ConnectMaster() error {
	// =========================================================================
	// ⚠️ TESTE TEMPORÁRIO: user e senha fixos no código, sem .env
	// Depois de testar, reverta para os.Getenv("MYSQLUSER") e
	// os.Getenv("MYSQLPASSWORD") e apague a senha em texto puro daqui.
	// =========================================================================
	user := "salao_app"
	pass := "Rk!4568_Salao@2026"

	host := os.Getenv("MYSQLHOST")
	port := os.Getenv("MYSQLPORT")
	dbName := os.Getenv("MYSQLDATABASE")

	// Logs de diagnóstico — não exibem a senha
	log.Printf("🔧 MYSQLUSER: %s", user)
	log.Printf("🔧 MYSQLHOST: %s", host)
	log.Printf("🔧 MYSQLPORT: %s", port)
	log.Printf("🔧 MYSQLDATABASE: %s", dbName)
	log.Printf("🔐 MYSQLPASSWORD configurada: %t", pass != "")
	log.Printf("🔐 Senha usada no teste tem %d caracteres", len(pass))

	if user == "" || pass == "" || host == "" || port == "" || dbName == "" {
		return fmt.Errorf("variáveis de ambiente do banco não configuradas")
	}

	// =========================================================================
	// CONFIGURAÇÃO DO MYSQL
	// =========================================================================

	// Usamos mysql.Config em vez de montar a string manualmente com fmt.Sprintf,
	// porque o driver faz o escaping correto de caracteres especiais na senha
	// (ex: "@"), evitando erro de parsing do DSN.
	cfg := mysql.Config{
		User:      user,
		Passwd:    pass,
		Net:       "tcp",
		Addr:      host + ":" + port,
		DBName:    dbName,
		ParseTime: true,
	}

	dsn := cfg.FormatDSN()

	log.Printf(
		"🔍 Conexão MySQL: %s:***@%s/%s",
		user,
		host+":"+port,
		dbName,
	)

	// =========================================================================
	// ABRIR CONEXÃO
	// =========================================================================

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com MySQL: %w", err)
	}

	// =========================================================================
	// TESTAR CONEXÃO
	// =========================================================================

	if err := db.Ping(); err != nil {
		db.Close()
		return fmt.Errorf("erro ao pingar banco MySQL: %w", err)
	}

	// =========================================================================
	// CONFIGURAÇÃO DO POOL
	// =========================================================================

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	MasterDB = db

	log.Println("===================================================")
	log.Println("✅ CONECTADO AO MYSQL COM SUCESSO")
	log.Println("===================================================")
	log.Println("📦 Banco:", dbName)
	log.Println("📡 Host:", host+":"+port)
	log.Println("👤 Usuário:", user)

	return nil
}

// CloseDB fecha a conexão com o banco
func CloseDB() error {
	if MasterDB != nil {
		return MasterDB.Close()
	}

	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// GetProjectByID retorna projeto pelo ID
func GetProjectByID(projectID int) (*Project, error) {
	var project Project

	query := "SELECT id, name, code, api_key FROM projects WHERE id = ? LIMIT 1"

	row := MasterDB.QueryRow(query, projectID)

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Code,
		&project.ApiKey,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProjectByApiKey retorna projeto pelo apiKey
func GetProjectByApiKey(apiKey string) (*Project, error) {
	var project Project

	query := "SELECT id, name, code, api_key FROM projects WHERE api_key = ? LIMIT 1"

	row := MasterDB.QueryRow(query, apiKey)

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Code,
		&project.ApiKey,
	)

	if err != nil {
		return nil, err
	}

	return &project, nil
}

// GetProjectCodeByID retorna o code de um projeto dado seu ID
func GetProjectCodeByID(projectID int) (string, error) {
	var code string

	query := "SELECT code FROM projects WHERE id = ? LIMIT 1"

	row := MasterDB.QueryRow(query, projectID)

	err := row.Scan(&code)

	if err != nil {
		return "", fmt.Errorf(
			"erro ao buscar code do projeto: %w",
			err,
		)
	}

	return code, nil
}

// BuildTableName constrói o nome completo da tabela com prefixo do projeto
func BuildTableName(project *Project, table string) string {
	return fmt.Sprintf("%s_%s", project.Code, table)
}

// RowsToMap converte sql.Rows para []map[string]interface{}
func RowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})

		for i, colName := range cols {
			val := columns[i]

			if b, ok := val.([]byte); ok {
				m[colName] = string(b)
			} else {
				m[colName] = val
			}
		}

		results = append(results, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// ============================================================================
// PROJECT STRUCT
// ============================================================================

type Project struct {
	ID     int
	Name   string
	Code   string
	ApiKey string
}