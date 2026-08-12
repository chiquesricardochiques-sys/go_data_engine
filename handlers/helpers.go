package handlers

import (
	"encoding/json"
	"net/http"
)

// ============================================================================
// HELPER FUNCTIONS - Funções auxiliares para handlers HTTP
// ============================================================================

// RespondSuccess envia resposta de sucesso em JSON
func RespondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// RespondError envia resposta de erro em JSON
func RespondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// RespondCreated envia resposta de criação bem-sucedida (201)
func RespondCreated(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

// RespondNoContent envia resposta sem conteúdo (204)
func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}
