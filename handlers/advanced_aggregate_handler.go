package handlers

import (
	"encoding/json"
	"net/http"
	"meu-provedor/models"
	"meu-provedor/services/data_service"
)

// ============================================================================
// AGGREGATE HANDLER
// ============================================================================

// AggregateHandler processa requisições de agregação (COUNT, SUM, AVG, MIN, MAX, EXISTS)
func AggregateHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AggregateRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Executar agregação
	result, err := services.ExecuteAggregate(req)
	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"result":  result,
	})
}