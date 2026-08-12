package handlers

import (
	"encoding/json"
	"net/http"

	"meu-provedor/services/data_service"
	"meu-provedor/models"
)

// ============================================================================
// SELECT HANDLERS
// ============================================================================




// AdvancedSelectHandler processa requisições de SELECT avançado
func AdvancedSelectHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdvancedSelectRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Executar SELECT
	result, err := services.ExecuteAdvancedSelect(req)
	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"data":    result,
		"count":   len(result),
	})
}

// AdvancedJoinSelectHandler processa requisições de SELECT com múltiplos JOINs
func AdvancedJoinSelectHandler(w http.ResponseWriter, r *http.Request) {
	var req models.AdvancedJoinSelectRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Executar JOIN SELECT
	result, err := services.ExecuteAdvancedJoinSelect(req)
	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"data":    result,
		"count":   len(result),
	})
}