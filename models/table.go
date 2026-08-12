package models

type ColumnRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
	Unique   bool   `json:"unique"`
}

type IndexRequest struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Type    string   `json:"type"` // UNIQUE ou INDEX
}

// CreateTableRequest - Agora usa project_id ao invés de project_code
type CreateTableRequest struct {
	ProjectID int64           `json:"project_id"`
	TableName string          `json:"table_name"`
	Columns   []ColumnRequest `json:"columns"`
	Indexes   []IndexRequest  `json:"indexes,omitempty"`
}

type ColumnDetail struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Nullable bool        `json:"nullable"`
	Default  interface{} `json:"default"`
	Key      string      `json:"key"`
	Extra    string      `json:"extra"`
}

type IndexDetail struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Type    string   `json:"type"`
}

type TableDetail struct {
	Name    string         `json:"name"`
	Columns []ColumnDetail `json:"columns"`
	Indexes []IndexDetail  `json:"indexes"`
}
