package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"meu-provedor/models"
	"meu-provedor/services/data_service"
)

// ============================================================================
// DELETE HANDLER
// ============================================================================

// DeleteHandler processa requisições de DELETE (hard ou soft)
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Normalizar modo (padrão: hard)
	mode := strings.ToLower(strings.TrimSpace(req.Mode))
	if mode == "" {
		mode = "hard"
	}

	var count int64
	var err error

	// Executar DELETE de acordo com o modo
	switch mode {
	case "soft":
		count, err = services.ExecuteSoftDelete(req)
	case "hard":
		count, err = services.ExecuteHardDelete(req)
	default:
		RespondError(w, "Modo inválido. Use 'soft' ou 'hard'", http.StatusBadRequest)
		return
	}

	if err != nil {
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar resultado
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Delete concluído",
		"mode":    mode,
		"count":   count,
	})
}