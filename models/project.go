package models

import "time"

// ============================================================================
// ENTITY MODELS - Estruturas de entidades do domínio
// ============================================================================

// Project - Representa um projeto no sistema
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`         // Prefixo único para tabelas (IMUTÁVEL)
	Description string    `json:"description"`
	ApiKey      string    `json:"api_key"`
	Type        string    `json:"type"`
	Version     string    `json:"version"`
	Status      string    `json:"status"`       // active, inactive, blocked
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectRequest - Request para criação de projeto
type ProjectRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`      // Obrigatório apenas na criação
	ApiKey  string `json:"api_key"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

// ProjectUpdateRequest - Request para atualização de projeto (SEM code)
type ProjectUpdateRequest struct {
	Name    string `json:"name"`
	ApiKey  string `json:"api_key"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Status  string `json:"status"`
}
