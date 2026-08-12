package models
import "errors"

// Column representa uma coluna com seu valor
type Column struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// InsertRequest - Requisição refatorada com especificação clara
type InsertRequest struct {
	ProjectID  int64    `json:"project_id"`
	InstanceID int64    `json:"id_instancia"`
	Table      string   `json:"table"`
	Columns    []Column `json:"columns"` // ✅ Agora é explícito: nome + valor
}

// BatchInsertRequest - Múltiplos inserts com mesma estrutura
type BatchInsertRequest struct {
	ProjectID  int64      `json:"project_id"`
	InstanceID int64      `json:"id_instancia"`
	Table      string     `json:"table"`
	Rows       [][]Column `json:"rows"` // ✅ Array de rows, cada row tem suas colunas
}

// Validate valida InsertRequest
func (r *InsertRequest) Validate() error {
	if r.ProjectID <= 0 {
		return errors.New("project_id inválido")
	}
	if r.InstanceID <= 0 {
		return errors.New("id_instancia inválido")
	}
	if r.Table == "" {
		return errors.New("table é obrigatória")
	}
	if len(r.Columns) == 0 {
		return errors.New("nenhuma coluna fornecida")
	}
	
	// Validar nomes das colunas
	for _, col := range r.Columns {
		if col.Name == "" {
			return errors.New("coluna com nome vazio")
		}
		if !IsValidColumnName(col.Name) {
			return errors.New("nome de coluna inválido: " + col.Name)
		}
	}
	
	return nil
}

// Validate valida BatchInsertRequest
func (r *BatchInsertRequest) Validate() error {
	if r.ProjectID <= 0 {
		return errors.New("project_id inválido")
	}
	if r.InstanceID <= 0 {
		return errors.New("id_instancia inválido")
	}
	if r.Table == "" {
		return errors.New("table é obrigatória")
	}
	if len(r.Rows) == 0 {
		return errors.New("nenhuma linha fornecida")
	}
	
	// Validar estrutura de cada row
	firstRowLen := len(r.Rows[0])
	for i, row := range r.Rows {
		if len(row) == 0 {
			return errors.New("linha vazia")
		}
		if len(row) != firstRowLen {
			return errors.New("linhas com número diferente de colunas")
		}
		
		// Validar cada coluna
		for _, col := range row {
			if col.Name == "" {
				return errors.New("coluna com nome vazio na linha " + string(rune(i)))
			}
			if !IsValidColumnName(col.Name) {
				return errors.New("nome de coluna inválido: " + col.Name)
			}
		}
	}
	
	return nil
}

// IsValidColumnName valida nome de coluna
func IsValidColumnName(name string) bool {
	if name == "" || len(name) > 64 {
		return false
	}
	
	// Aceita: letras, números, underscore
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || 
			(char >= 'A' && char <= 'Z') || 
			(char >= '0' && char <= '9') || 
			char == '_') {
			return false
		}
	}
	
	return true
}


// ============================================================================
// REQUEST MODELS - Estruturas de requisição HTTP
// ============================================================================

// DeleteRequest - Requisição para deletar registros
type DeleteRequest struct {
	ProjectID  int64                  `json:"project_id"`
	InstanceID int64                  `json:"id_instancia"`
	Table      string                 `json:"table"`
	Where      map[string]interface{} `json:"where,omitempty"`
	WhereRaw   string                 `json:"where_raw,omitempty"`
	Mode       string                 `json:"mode,omitempty"` // "hard" ou "soft"
}

// AdvancedSelectRequest - Requisição para SELECT avançado
type AdvancedSelectRequest struct {
	ProjectID  int64                  `json:"project_id"`
	InstanceID int64                  `json:"id_instancia"`
	Table      string                 `json:"table"`
	Alias      string                 `json:"alias,omitempty"`
	Select     []string               `json:"select,omitempty"`
	Joins      []Join                 `json:"joins,omitempty"`
	Where      map[string]interface{} `json:"where,omitempty"`
	WhereRaw   string                 `json:"where_raw,omitempty"`
	GroupBy    string                 `json:"group_by,omitempty"`
	Having     string                 `json:"having,omitempty"`
	OrderBy    string                 `json:"order_by,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Offset     int                    `json:"offset,omitempty"`
}

// Join - Configuração de JOIN
type Join struct {
	Type  string `json:"type"`  // INNER, LEFT, RIGHT
	Table string `json:"table"`
	Alias string `json:"alias,omitempty"`
	On    string `json:"on"`
}

// AdvancedJoinSelectRequest - Requisição para SELECT com múltiplos JOINs
type AdvancedJoinSelectRequest struct {
	ProjectID  int64                  `json:"project_id"`
	InstanceID int64                  `json:"id_instancia"`
	Base       JoinBase               `json:"base"`
	Joins      []JoinItem             `json:"joins,omitempty"`
	Where      map[string]interface{} `json:"where,omitempty"`
	WhereRaw   []string               `json:"where_raw,omitempty"`
	GroupBy    string                 `json:"group_by,omitempty"`
	Having     string                 `json:"having,omitempty"`
	OrderBy    string                 `json:"order_by,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Offset     int                    `json:"offset,omitempty"`
}

// JoinBase - Tabela base para JOIN
type JoinBase struct {
	Table   string   `json:"table"`
	Alias   string   `json:"alias,omitempty"`
	Columns []string `json:"columns,omitempty"`
}

// JoinItem - Item de JOIN
type JoinItem struct {
	Type    string   `json:"type"`
	Table   string   `json:"table"`
	Alias   string   `json:"alias,omitempty"`
	On      string   `json:"on"`
	Columns []string `json:"columns,omitempty"`
}



// UpdateRequest - Requisição para UPDATE
type UpdateRequest struct {
	ProjectID  int64                  `json:"project_id"`
	InstanceID int64                  `json:"id_instancia"`
	Table      string                 `json:"table"`
	Data       map[string]interface{} `json:"data"`
	Where      map[string]interface{} `json:"where,omitempty"`
	WhereRaw   string                 `json:"where_raw,omitempty"`
}

// BatchUpdateRequest - Requisição para UPDATE em lote
type BatchUpdateRequest struct {
	ProjectID  int64               `json:"project_id"`
	InstanceID int64               `json:"id_instancia"`
	Table      string              `json:"table"`
	Updates    []UpdateItem        `json:"updates"`
}

// UpdateItem - Item individual de update em lote
type UpdateItem struct {
	Data  map[string]interface{} `json:"data"`
	Where map[string]interface{} `json:"where"`
}

// AggregateRequest - Requisição para operações de agregação
type AggregateRequest struct {
	ProjectID  int64                  `json:"project_id"`
	InstanceID int64                  `json:"id_instancia"`
	Table      string                 `json:"table"`
	Operation  string                 `json:"operation"` // COUNT, SUM, AVG, MIN, MAX, EXISTS
	Column     string                 `json:"column,omitempty"`
	Where      map[string]interface{} `json:"where,omitempty"`
}

// ============================================================================
// VALIDATION METHODS
// ============================================================================

// Validate - Valida DeleteRequest
func (r *DeleteRequest) Validate() error {
	if r.ProjectID <= 0 {
		return ErrInvalidProjectID
	}
	if r.InstanceID <= 0 {
		return ErrInvalidInstanceID
	}
	if r.Table == "" {
		return ErrTableRequired
	}
	return nil
}

// Validate - Valida AdvancedSelectRequest
func (r *AdvancedSelectRequest) Validate() error {
	if r.ProjectID <= 0 {
		return ErrInvalidProjectID
	}
	if r.InstanceID <= 0 {
		return ErrInvalidInstanceID
	}
	if r.Table == "" {
		return ErrTableRequired
	}
	return nil
}



// Validate - Valida UpdateRequest
func (r *UpdateRequest) Validate() error {
	if r.ProjectID <= 0 {
		return ErrInvalidProjectID
	}
	if r.InstanceID <= 0 {
		return ErrInvalidInstanceID
	}
	if r.Table == "" {
		return ErrTableRequired
	}
	if len(r.Data) == 0 {
		return ErrNoDataProvided
	}
	return nil
}

// Validate - Valida AggregateRequest
func (r *AggregateRequest) Validate() error {
	if r.ProjectID <= 0 {
		return ErrInvalidProjectID
	}
	if r.InstanceID <= 0 {
		return ErrInvalidInstanceID
	}
	if r.Table == "" {
		return ErrTableRequired
	}
	if r.Operation == "" {
		return ErrOperationRequired
	}
	return nil

}
