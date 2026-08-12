package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"meu-provedor/config"
	"meu-provedor/routes"
)

// ============================================================================
// MAIN APPLICATION
// ============================================================================

func main() {
	log.Println("═══════════════════════════════════════════════════════════")
	log.Println("  SISTEMA DE GERENCIAMENTO MULTI-PROJETO")
	log.Println("═══════════════════════════════════════════════════════════")

	// 1️⃣ Carregar variáveis de ambiente
	config.LoadEnv()

	// 2️⃣ Conectar ao banco de dados
	if err := config.ConnectMaster(); err != nil {
		log.Fatalf("❌ Falha ao conectar ao banco: %v", err)
	}
	defer config.CloseDB()

	// 3️⃣ Definir porta do servidor
	port := config.GetEnvOrDefault("PORT", "8080")

	// 4️⃣ Configurar graceful shutdown
	go handleShutdown()

	// 5️⃣ Iniciar servidor HTTP
	routes.StartServer(port)
}

// handleShutdown gerencia o desligamento gracioso do servidor
func handleShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("\n⚠️ Sinal de shutdown recebido")
	
	// Fechar conexão com banco
	if err := config.CloseDB(); err != nil {
		log.Printf("❌ Erro ao fechar banco: %v", err)
	} else {
		log.Println("✅ Banco de dados desconectado")
	}

	log.Println("👋 Servidor encerrado")
	os.Exit(0)
}