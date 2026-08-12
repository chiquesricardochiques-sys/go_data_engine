package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"meu-provedor/handlers"
	"meu-provedor/security"
)

// ============================================================================
// ROUTER SETUP
// ============================================================================

// SetupRouter configura todas as rotas da API
func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Aplicar CORS globalmente
	r.Use(security.CORS)

	// Criar subrouter protegido
	protected := r.PathPrefix("/").Subrouter()
	protected.Use(security.InternalOnly)

	// ========================================
	// DATA ENGINE ROUTES
	// ========================================

	// SELECT InsertDebugHandler
	protected.HandleFunc("/data/select", handlers.AdvancedSelectHandler).Methods("POST")

	protected.HandleFunc("/data/join-select", handlers.AdvancedJoinSelectHandler).Methods("POST")

	// INSERT
	protected.HandleFunc("/data/insert", handlers.InsertHandler).Methods("POST")
	protected.HandleFunc("/data/batch-insert", handlers.BatchInsertHandler).Methods("POST")

	// UPDATE
	protected.HandleFunc("/data/update", handlers.UpdateHandler).Methods("POST")
	protected.HandleFunc("/data/batch-update", handlers.BatchUpdateHandler).Methods("POST")

	// DELETE
	protected.HandleFunc("/data/delete", handlers.DeleteHandler).Methods("POST")

	// AGGREGATE
	protected.HandleFunc("/data/aggregate", handlers.AggregateHandler).Methods("POST")


	/*
	====================================================
	PROJETOS
	====================================================
	*/

	protected.HandleFunc("/projects", handlers.ListProjects).Methods("GET")
	protected.HandleFunc("/projects", handlers.CreateProject).Methods("POST")
	protected.HandleFunc("/projects/{id}", handlers.UpdateProject).Methods("PUT")
	protected.HandleFunc("/projects/{id}", handlers.DeleteProject).Methods("DELETE")

	/*
	====================================================
	INSTÂNCIAS
	====================================================
	*/

	protected.HandleFunc("/instances", handlers.ListInstances).Methods("GET")
	protected.HandleFunc("/instances", handlers.CreateInstance).Methods("POST")
	protected.HandleFunc("/instances/{id}", handlers.UpdateInstance).Methods("PUT")
	protected.HandleFunc("/instances/{id}", handlers.DeleteInstance).Methods("DELETE")

	/*
	====================================================
	SCHEMA – TABELAS
	====================================================
	*/

	protected.HandleFunc("/schema/table", handlers.CreateProjectTable).Methods("POST")
	protected.HandleFunc("/schema/tables", handlers.ListProjectTables).Methods("GET")
	protected.HandleFunc("/schema/table/details", handlers.GetTableDetails).Methods("GET")
	protected.HandleFunc("/schema/table", handlers.DeleteProjectTable).Methods("DELETE")

	/*
	====================================================
	SCHEMA – COLUNAS
	====================================================
	*/

	protected.HandleFunc("/schema/column", handlers.AddColumn).Methods("POST")
	protected.HandleFunc("/schema/column", handlers.ModifyColumn).Methods("PUT")
	protected.HandleFunc("/schema/column", handlers.DropColumn).Methods("DELETE")

	/*
	====================================================
	SCHEMA – ÍNDICES
	====================================================
	*/

	protected.HandleFunc("/schema/index", handlers.AddIndex).Methods("POST")
	protected.HandleFunc("/schema/index", handlers.DropIndex).Methods("DELETE")

	return r
}

func StartServer(port string) {
	r := SetupRouter()
	log.Println("🚀 Servidor iniciado na porta", port)

	if err := http.ListenAndServe("0.0.0.0:"+port, r); err != nil {
		log.Fatal("❌ Erro ao iniciar servidor:", err)
	}
}



