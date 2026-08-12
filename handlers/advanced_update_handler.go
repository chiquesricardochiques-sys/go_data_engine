package handlers

import (
	"encoding/json"
	"net/http"
	"meu-provedor/models"
	"meu-provedor/services/data_service"
)

// ============================================================================
// UPDATE HANDLERS
// ============================================================================

// UpdateHandler processa requisições de UPDATE
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Executar UPDATE
	count, err := services.ExecuteUpdate(req)
	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Atualização concluída",
		"count":   count,
	})
}

// BatchUpdateHandler processa requisições de UPDATE em lote
func BatchUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var req models.BatchUpdateRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Executar BATCH UPDATE
	count, err := services.ExecuteBatchUpdate(req)
	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Atualizações concluídas",
		"count":   count,
	})
}