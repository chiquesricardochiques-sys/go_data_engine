package services

import (
	"database/sql"
	"fmt"
	"meu-provedor/config"
	"meu-provedor/engine/query"
	"meu-provedor/models"
)

// ============================================================================
// AGGREGATE SERVICE
// ============================================================================

// ExecuteAggregate executa operações de agregação (COUNT, SUM, AVG, MIN, MAX, EXISTS)
func ExecuteAggregate(req models.AggregateRequest) (interface{}, error) {
	// Validar requisição
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Obter código do projeto
	projectCode, err := GetProjectCodeByID(req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Construir nome da tabela
	table, err := BuildTableName(projectCode, req.Table)
	if err != nil {
		return nil, err
	}

	// Criar AggregateBuilder
	builder := query.NewAggregate(table, "", req.Operation, req.Column)

	// Filtro obrigatório: id_instancia
	builder.AddWhere("id_instancia = ?", req.InstanceID)

	// Adicionar filtros simples
	for col, val := range req.Where {
		if !query.IsValidColumnName(col) {
			return nil, fmt.Errorf("%w: %s", models.ErrInvalidColumn, col)
		}
		builder.AddWhere(col+" = ?", val)
	}

	// Executar query
	sqlQuery := builder.Build()
	var result interface{}
	
	err = config.MasterDB.QueryRow(sqlQuery, builder.GetValues()...).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoResultsFound
		}
		return nil, fmt.Errorf("%w: %v", models.ErrQueryFailed, err)
	}

	if result == nil {
		return nil, models.ErrNoResultsFound
	}

	return result, nil
}