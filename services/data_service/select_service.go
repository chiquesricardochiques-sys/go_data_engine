package services

import (
	"fmt"
	"meu-provedor/config"
	"meu-provedor/engine/query"
	"meu-provedor/models"
)

// ============================================================================
// SELECT SERVICE
// ============================================================================

// ExecuteAdvancedSelect executa um SELECT avançado com suporte a JOINs
func ExecuteAdvancedSelect(req models.AdvancedSelectRequest) ([]map[string]interface{}, error) {
	// Validar requisição
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Regra invertida:
	// Se req.UsePrefix for igual a 1, ignoramos o prefixo (usamos a tabela pura).
	// Se for nil ou 0, aplicamos o BuildTableName normalmente (com prefixo).
	ignorarPrefixo := false
	if req.UsePrefix != nil && *req.UsePrefix == 1 {
		ignorarPrefixo = true
	}

	// Variável para a tabela principal
	var mainTable string
	var err error

	if ignorarPrefixo {
		mainTable = req.Table
	} else {
		// Obter código do projeto
		projectCode, err := GetProjectCodeByID(req.ProjectID)
		if err != nil {
			return nil, err
		}

		// Construir nome físico da tabela com prefixo
		mainTable, err = BuildTableName(projectCode, req.Table)
		if err != nil {
			return nil, err
		}
	}

	mainAlias := req.Alias
	if mainAlias == "" {
		mainAlias = mainTable
	}

	// Criar SelectBuilder
	builder := query.NewSelect(mainTable, mainAlias)

	// Definir colunas
	if len(req.Select) > 0 {
		builder.SetColumns(req.Select)
	}

	// Adicionar JOINs respeitando a mesma regra
	for _, j := range req.Joins {
		var joinTable string
		if ignorarPrefixo {
			joinTable = j.Table
		} else {
			projectCode, err := GetProjectCodeByID(req.ProjectID)
			if err != nil {
				return nil, err
			}
			joinTable, err = BuildTableName(projectCode, j.Table)
			if err != nil {
				return nil, err
			}
		}
		builder.AddJoin(j.Type, joinTable, j.Alias, j.On)
	}

	// Filtro obrigatório: id_instancia (só adiciona se NÃO for para ignorar o prefixo, ou conforme sua regra de negócio)
	if !(req.UsePrefix != nil && *req.UsePrefix == 1) {
		builder.AddWhere(fmt.Sprintf("%s.id_instancia = ?", mainAlias), req.InstanceID)
	}

	// Filtros simples (WHERE)
	for k, v := range req.Where {
		if !query.IsValidColumnName(k) {
			return nil, fmt.Errorf("%w: %s", models.ErrInvalidColumn, k)
		}
		builder.AddWhere(k+" = ?", v)
	}

	// Filtro raw (WHERE customizado)
	if req.WhereRaw != "" {
		builder.AddWhere("(" + req.WhereRaw + ")")
	}

	// GROUP BY
	if req.GroupBy != "" {
		builder.SetGroupBy(req.GroupBy)
	}

	// HAVING
	if req.Having != "" {
		builder.SetHaving(req.Having)
	}

	// ORDER BY
	if req.OrderBy != "" {
		builder.SetOrderBy(req.OrderBy)
	}

	// LIMIT e OFFSET
	if req.Limit > 0 {
		builder.SetLimitOffset(req.Limit, req.Offset)
	}

	sql := builder.Build()
	values := builder.GetValues()

	// Executar query
	rows, err := config.MasterDB.Query(sql, values...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", models.ErrQueryFailed, err)
	}
	defer rows.Close()

	// Converter para map
	result, err := RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}
