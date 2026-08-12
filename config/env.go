package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// ============================================================================
// ENVIRONMENT CONFIGURATION
// ============================================================================

// LoadEnv carrega variáveis de ambiente do arquivo .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ Arquivo .env não encontrado, usando variáveis do sistema")
	} else {
		log.Println("✅ Variáveis de ambiente carregadas")
	}
}

// GetEnv retorna valor de variável de ambiente obrigatória
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ ERRO: Variável de ambiente %s não definida", key)
	}
	return value
}

// GetEnvOrDefault retorna valor de variável ou default se não existir
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}