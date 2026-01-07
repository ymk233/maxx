package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Bowl42/maxx-next/internal/adapter/client"
	_ "github.com/Bowl42/maxx-next/internal/adapter/provider/antigravity" // Register antigravity adapter
	_ "github.com/Bowl42/maxx-next/internal/adapter/provider/custom"      // Register custom adapter
	"github.com/Bowl42/maxx-next/internal/executor"
	"github.com/Bowl42/maxx-next/internal/handler"
	"github.com/Bowl42/maxx-next/internal/repository/cached"
	"github.com/Bowl42/maxx-next/internal/repository/sqlite"
	"github.com/Bowl42/maxx-next/internal/router"
)

func main() {
	// Parse flags
	addr := flag.String("addr", ":8080", "Server address")
	dbPath := flag.String("db", "maxx.db", "SQLite database path")
	flag.Parse()

	// Initialize database
	db, err := sqlite.NewDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create repositories
	providerRepo := sqlite.NewProviderRepository(db)
	routeRepo := sqlite.NewRouteRepository(db)
	projectRepo := sqlite.NewProjectRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)
	retryConfigRepo := sqlite.NewRetryConfigRepository(db)
	routingStrategyRepo := sqlite.NewRoutingStrategyRepository(db)
	proxyRequestRepo := sqlite.NewProxyRequestRepository(db)
	attemptRepo := sqlite.NewProxyUpstreamAttemptRepository(db)

	// Create cached repositories
	cachedProviderRepo := cached.NewProviderRepository(providerRepo)
	cachedRouteRepo := cached.NewRouteRepository(routeRepo)
	cachedRetryConfigRepo := cached.NewRetryConfigRepository(retryConfigRepo)
	cachedRoutingStrategyRepo := cached.NewRoutingStrategyRepository(routingStrategyRepo)
	cachedSessionRepo := cached.NewSessionRepository(sessionRepo)

	// Load cached data
	if err := cachedProviderRepo.Load(); err != nil {
		log.Printf("Warning: Failed to load providers cache: %v", err)
	}
	if err := cachedRouteRepo.Load(); err != nil {
		log.Printf("Warning: Failed to load routes cache: %v", err)
	}
	if err := cachedRetryConfigRepo.Load(); err != nil {
		log.Printf("Warning: Failed to load retry configs cache: %v", err)
	}
	if err := cachedRoutingStrategyRepo.Load(); err != nil {
		log.Printf("Warning: Failed to load routing strategies cache: %v", err)
	}

	// Create router
	r := router.NewRouter(cachedRouteRepo, cachedProviderRepo, cachedRoutingStrategyRepo, cachedRetryConfigRepo)

	// Create executor
	exec := executor.NewExecutor(r, proxyRequestRepo, attemptRepo, cachedRetryConfigRepo)

	// Create client adapter
	clientAdapter := client.NewAdapter()

	// Create handlers
	proxyHandler := handler.NewProxyHandler(clientAdapter, exec, cachedSessionRepo)
	adminHandler := handler.NewAdminHandler(
		cachedProviderRepo,
		cachedRouteRepo,
		projectRepo,
		cachedSessionRepo,
		cachedRetryConfigRepo,
		cachedRoutingStrategyRepo,
		proxyRequestRepo,
	)

	// Setup routes
	mux := http.NewServeMux()

	// Admin API routes
	mux.Handle("/admin/", adminHandler)

	// Proxy routes - catch all AI API endpoints
	// Claude API
	mux.Handle("/v1/messages", proxyHandler)
	// OpenAI API
	mux.Handle("/v1/chat/completions", proxyHandler)
	// Codex API
	mux.Handle("/v1/responses", proxyHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	log.Printf("Starting maxx-next server on %s", *addr)
	log.Printf("Admin API: http://localhost%s/admin/", *addr)
	log.Printf("Proxy endpoints:")
	log.Printf("  Claude: http://localhost%s/v1/messages", *addr)
	log.Printf("  OpenAI: http://localhost%s/v1/chat/completions", *addr)
	log.Printf("  Codex:  http://localhost%s/v1/responses", *addr)

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
