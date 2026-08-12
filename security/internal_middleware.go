package security

import (
	"net/http"
	"os"
)

// ============================================================================
// SECURITY MIDDLEWARE
// ============================================================================

// InternalOnly valida token interno antes de processar requisição
func InternalOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		internalToken := os.Getenv("INTERNAL_TOKEN")
		requestToken := r.Header.Get("X-Internal-Token")

		// Validar token
		if requestToken == "" || requestToken != internalToken {
			http.Error(w, "Acesso negado - token inválido ou ausente", http.StatusForbidden)
			return
		}

		// Token válido, prosseguir
		next.ServeHTTP(w, r)
	})
}

// CORS adiciona headers CORS à resposta
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Internal-Token")

		// Tratar preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}