package security

import (
    "errors"
    "log"
    "meu-provedor/config"
)

// Valida a API KEY recebida e retorna o projeto correspondente
func ValidateApiKey(apiKey string) (*config.Project, error) {
    if apiKey == "" {
        return nil, errors.New("⚠️ API KEY não fornecida")
    }

    project, err := config.GetProjectByApiKey(apiKey)
    if err != nil {
        log.Println("❌ API KEY inválida:", apiKey)
        return nil, errors.New("API KEY inválida")
    }

    return project, nil
}
