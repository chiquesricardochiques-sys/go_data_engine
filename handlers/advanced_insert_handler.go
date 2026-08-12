package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"meu-provedor/models"
	"meu-provedor/services/data_service"
)

// InsertHandler processa INSERT único
func InsertHandler(w http.ResponseWriter, r *http.Request) {
	var req models.InsertRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Erro ao decodificar JSON: %v", err)
		RespondError(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Executar INSERT
	lastID, err := services.ExecuteInsert(req)
	if err != nil {
		log.Printf("❌ Erro ao executar INSERT: %v", err)
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Retornar sucesso
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Registro inserido com sucesso",
		"id":      lastID,
	})
}

// BatchInsertHandler processa múltiplos INSERTs
func BatchInsertHandler(w http.ResponseWriter, r *http.Request) {
	var req models.BatchInsertRequest
	
	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Erro ao decodificar JSON: %v", err)
		RespondError(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Executar BATCH INSERT
	count, err := services.ExecuteBatchInsert(req)
	if err != nil {
		log.Printf("❌ Erro ao executar BATCH INSERT: %v", err)
		RespondError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Retornar sucesso
	RespondSuccess(w, map[string]interface{}{
		"success": true,
		"message": "Registros inseridos com sucesso",
		"count":   count,
	})
}


