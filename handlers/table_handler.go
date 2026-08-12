package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"meu-provedor/models"
	tableService "meu-provedor/services/table"
)

func CreateProjectTable(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if req.ProjectID <= 0 {
		http.Error(w, "project_id is required", 400)
		return
	}

	tableName, err := tableService.Create(req.ProjectID, req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "TABLE CREATED",
		"table":   tableName,
	})
}

func ListProjectTables(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		http.Error(w, "project_id is required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	tables, err := tableService.List(projectID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(tables)
}

func DeleteProjectTable(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")

	if projectIDStr == "" || tableName == "" {
		http.Error(w, "project_id and table are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	if err := tableService.Delete(projectID, tableName); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("TABLE DELETED"))
}

// GET TABLE DETAILS
func GetTableDetails(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")

	if projectIDStr == "" || tableName == "" {
		http.Error(w, "project_id and table are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	details, err := tableService.GetDetails(projectID, tableName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(details)
}

// ADD COLUMN
func AddColumn(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")

	if projectIDStr == "" || tableName == "" {
		http.Error(w, "project_id and table are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	var col tableService.ColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if err := tableService.AddColumn(projectID, tableName, col); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("COLUMN ADDED"))
}

// MODIFY COLUMN
func ModifyColumn(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")

	if projectIDStr == "" || tableName == "" {
		http.Error(w, "project_id and table are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	var col tableService.ColumnRequest
	if err := json.NewDecoder(r.Body).Decode(&col); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if err := tableService.ModifyColumn(projectID, tableName, col); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("COLUMN MODIFIED"))
}

// DROP COLUMN
func DropColumn(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")
	columnName := r.URL.Query().Get("column")

	if projectIDStr == "" || tableName == "" || columnName == "" {
		http.Error(w, "project_id, table and column are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	if err := tableService.DropColumn(projectID, tableName, columnName); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("COLUMN DROPPED"))
}

// ADD INDEX
func AddIndex(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")

	if projectIDStr == "" || tableName == "" {
		http.Error(w, "project_id and table are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	var idx tableService.IndexRequest
	if err := json.NewDecoder(r.Body).Decode(&idx); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if err := tableService.AddIndex(projectID, tableName, idx); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("INDEX ADDED"))
}

// DROP INDEX
func DropIndex(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	tableName := r.URL.Query().Get("table")
	indexName := r.URL.Query().Get("index")

	if projectIDStr == "" || tableName == "" || indexName == "" {
		http.Error(w, "project_id, table and index are required", 400)
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", 400)
		return
	}

	if err := tableService.DropIndex(projectID, tableName, indexName); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("INDEX DROPPED"))
}
