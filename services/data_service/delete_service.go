package services

import (
	"fmt"
	"time"
	"meu-provedor/config"
	"meu-provedor/engine/query"
	"meu-provedor/models"
)

// ============================================================================
// DELETE SERVICE
// ============================================================================

// ExecuteHardDelete executa um DELETE físico (remove do banco)
func ExecuteHardDelete(req models.DeleteRequest) (int64, error) {
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

	// Criar DeleteBuilder
	builder := query.NewDelete(table)

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

	// Executar DELETE
	sqlQuery, args := builder.Build()
	result, err := config.MasterDB.Exec(sqlQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrDeleteFailed, err)
	}

	// Retornar quantidade de linhas afetadas
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ExecuteSoftDelete executa um soft delete (marca como deletado)
func ExecuteSoftDelete(req models.DeleteRequest) (int64, error) {
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

	// Garantir que a coluna deleted_at existe
	if err := EnsureSoftDeleteColumn(config.MasterDB, table); err != nil {
		return 0, err
	}

	// Criar SoftDeleteBuilder
	builder := query.NewSoftDelete(table)

	// Filtro obrigatório: id_instancia
	builder.AddWhere("id_instancia = ?", req.InstanceID)

	// Adicionar filtros simples
	for col, val := range req.Where {
		if !query.IsValidColumnName(col) {
			return 0, fmt.Errorf("%w: %s", models.ErrInvalidColumn, col)
		}
		builder.AddWhere(col+" = ?", val)
	}

	// Filtro raw opcional
	if req.WhereRaw != "" {
		builder.AddRawWhere(req.WhereRaw)
	}

	// Executar soft delete
	sqlQuery, args := builder.Build(time.Now())
	result, err := config.MasterDB.Exec(sqlQuery, args...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", models.ErrDeleteFailed, err)
	}

	// Retornar quantidade de linhas afetadas
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}