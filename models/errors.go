package models

import "errors"

// ============================================================================
// ERROR DEFINITIONS
// ============================================================================

var (
	// Erros de validação de requisição
	ErrInvalidProjectID   = errors.New("project_id inválido ou não informado")
	ErrInvalidInstanceID  = errors.New("id_instancia inválido ou não informado")
	ErrTableRequired      = errors.New("nome da tabela é obrigatório")
	ErrNoDataProvided     = errors.New("nenhum dado fornecido")
	ErrOperationRequired  = errors.New("operação de agregação é obrigatória")
	ErrInvalidColumn      = errors.New("nome de coluna inválido")
	ErrInvalidIdentifier  = errors.New("identificador contém caracteres inválidos")

	// Erros de projeto
	ErrProjectNotFound    = errors.New("projeto não encontrado")
	ErrInvalidAPIKey      = errors.New("API key inválida")
	ErrAPIKeyNotProvided  = errors.New("API key não fornecida")

	// Erros de operação
	ErrQueryFailed        = errors.New("falha ao executar query")
	ErrNoResultsFound     = errors.New("nenhum resultado encontrado")
	ErrInsertFailed       = errors.New("falha ao inserir dados")
	ErrUpdateFailed       = errors.New("falha ao atualizar dados")
	ErrDeleteFailed       = errors.New("falha ao deletar dados")

	// Erros de conexão
	ErrDatabaseConnection = errors.New("erro de conexão com banco de dados")
	// erro para projetos
	ErrInvalidProjectData = errors.New("dados do projeto inválidos")
	ErrProjectCodeExists  = errors.New("project code já existe")

)
