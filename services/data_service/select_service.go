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

	// Obter código do projeto
	projectCode, err := GetProjectCodeByID(req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Construir nome físico da tabela
	mainTable, err := BuildTableName(projectCode, req.Table)
	if err != nil {
		return nil, err
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

	// Adicionar JOINs
	for _, j := range req.Joins {
		joinTable, err := BuildTableName(projectCode, j.Table)
		if err != nil {
			return nil, err
		}
		builder.AddJoin(j.Type, joinTable, j.Alias, j.On)
	}

	// Filtro obrigatório: id_instancia
	builder.AddWhere(fmt.Sprintf("%s.id_instancia = ?", mainAlias), req.InstanceID)

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

	// ============================================================================
	// DEBUG SQL
	// ============================================================================

	sql := builder.Build()
	values := builder.GetValues()

	fmt.Println("======================================================")
	fmt.Println("🟢 SQL GERADA")
	fmt.Println(sql)
	fmt.Println("------------------------------------------------------")
	fmt.Printf("📦 Valores: %#v\n", values)
	fmt.Println("======================================================")

	// Executar query
	rows, err := config.MasterDB.Query(sql, values...)
	if err != nil {

		fmt.Println("======================================================")
		fmt.Println("🔴 ERRO AO EXECUTAR SQL")
		fmt.Println(sql)
		fmt.Println("------------------------------------------------------")
		fmt.Printf("📦 Valores: %#v\n", values)
		fmt.Println("------------------------------------------------------")
		fmt.Println(err)
		fmt.Println("======================================================")

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
