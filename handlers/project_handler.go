package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	projectService "meu-provedor/services/project"
	"meu-provedor/models"
)

func CreateProject(w http.ResponseWriter, r *http.Request) {
	var req models.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if err := projectService.Create(req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "PROJECT CREATED",
	})
}

func ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := projectService.List()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(projects)
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	var req models.ProjectUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}

	if err := projectService.Update(id, req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("PROJECT UPDATED"))
}

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)

	if err := projectService.Delete(id); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("PROJECT DELETED"))
}
