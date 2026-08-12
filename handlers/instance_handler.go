package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"meu-provedor/models"
	"meu-provedor/services/instance"
)

func CreateInstance(w http.ResponseWriter, r *http.Request) {
	var req models.InstanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := instance.Create(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("INSTANCE CREATED"))
}

func ListInstances(w http.ResponseWriter, r *http.Request) {
	var projectID *int64

	if pid := r.URL.Query().Get("project_id"); pid != "" {
		id, err := strconv.ParseInt(pid, 10, 64)
		if err != nil {
			http.Error(w, "invalid project_id", http.StatusBadRequest)
			return
		}
		projectID = &id
	}

	instances, err := instance.List(projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instances)
}

func UpdateInstance(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req models.InstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := instance.Update(id, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("INSTANCE UPDATED"))
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := instance.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("INSTANCE DELETED"))
}
