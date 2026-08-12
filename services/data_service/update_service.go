package services

import (
	"fmt"
	"meu-provedor/config"
	"meu-provedor/engine/query"
	"meu-provedor/models"
)

// ============================================================================
// UPDATE SERVICE
// ============================================================================

// ExecuteUpdate executa um UPDATE
func ExecuteUpdate(req models.UpdateRequest) (int64, error) {
	// Validar requisição
	if err := req.Validate(); err != nil {
		return 0, err
	}

	// Obter código do projeto
	projectCode, err := GetProjectCodeByID(req.ProjectID)
	if err != nil {
		return 0, err
	}

	// Construir nome da tabela
	table, err := BuildTableName(projectCode, req.Table)
	if err != nil {
		return 0, err
	}

	// Criar UpdateBuilder
	builder := query.NewUpdate(table)

	// Adicionar campos a atualizar
	for col, val := range req.Data {
		if !query.IsValidColumnName(col) {
			return 0, fmt.Errorf("%w: %s", models.ErrInvalidColumn, col)
		}
		builder.Set(col, val)
	}

	// Filtro obrigatório: id_instancia
	builder.Where("id_instancia = ?", req.InstanceID)

	// Adicionar filtros simples
	for col, val := range req.Where {
		if !query.IsValidColumnName(col) {
			return 0, fmt.Errorf("%w: %s", models.ErrInvalidColumn, col)
		}
		builder.Where(col+" = ?", val)
	}

	// Filtro raw opcional
	if req.WhereRaw != "" {
		builder.WhereRaw(req.WhereRaw)
	}

	// Executar UPDATE
	sqlQuery, args := builder.Build()
	result, err := config.MasterDB.Exec(sqlQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrUpdateFailed, err)
	}

	// Retornar quantidade de linhas afetadas
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ExecuteBatchUpdate executa múltiplos UPDATEs
func ExecuteBatchUpdate(req models.BatchUpdateRequest) (int64, error) {
	// Validar requisição básica
	if req.ProjectID <= 0 {
		return 0, models.ErrInvalidProjectID
	}
	if req.InstanceID <= 0 {
		return 0, models.ErrInvalidInstanceID
	}
	if req.Table == "" {
		return 0, models.ErrTableRequired
	}
	if len(req.Updates) == 0 {
		return 0, models.ErrNoDataProvided
	}

	// Obter código do projeto
	projectCode, err := GetProjectCodeByID(req.ProjectID)
	if err != nil {
		return 0, err
	}

	// Construir nome da tabela
	table, err := BuildTableName(projectCode, req.Table)
	if err != nil {
		return 0, err
	}

	var totalAffected int64

	// Executar cada update individualmente
	for _, update := range req.Updates {
		builder := query.NewUpdate(table)

		// Adicionar campos a atualizar
		for col, val := range update.Data {
			if !query.IsValidColumnName(col) {
				return 0, fmt.Errorf("%w: %s", models.ErrInvalidColumn, col)
			}
			builder.Set(col, val)
		}

		// Filtro obrigatório: id_instancia
		builder.Where("id_instancia = ?", req.InstanceID)

		// Adicionar filtros do update
		for col, val := range update.Where {
			if !query.IsValidColumnName(col) {
				return 0, fmt.Errorf("%w: %s", models.ErrInvalidColumn, col)
			}
			builder.Where(col+" = ?", val)
		}

		// Executar
		sqlQuery, args := builder.Build()
		result, err := config.MasterDB.Exec(sqlQuery, args...)
		if err != nil {
			return totalAffected, fmt.Errorf("%w: %v", models.ErrUpdateFailed, err)
		}

		affected, _ := result.RowsAffected()
		totalAffected += affected
	}

	return totalAffected, nil
}