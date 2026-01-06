package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bowl42/maxx-next/internal/domain"
	"github.com/Bowl42/maxx-next/internal/repository"
)

// AdminHandler handles admin API requests
type AdminHandler struct {
	providerRepo        repository.ProviderRepository
	routeRepo           repository.RouteRepository
	projectRepo         repository.ProjectRepository
	sessionRepo         repository.SessionRepository
	retryConfigRepo     repository.RetryConfigRepository
	routingStrategyRepo repository.RoutingStrategyRepository
	proxyRequestRepo    repository.ProxyRequestRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	providerRepo repository.ProviderRepository,
	routeRepo repository.RouteRepository,
	projectRepo repository.ProjectRepository,
	sessionRepo repository.SessionRepository,
	retryConfigRepo repository.RetryConfigRepository,
	routingStrategyRepo repository.RoutingStrategyRepository,
	proxyRequestRepo repository.ProxyRequestRepository,
) *AdminHandler {
	return &AdminHandler{
		providerRepo:        providerRepo,
		routeRepo:           routeRepo,
		projectRepo:         projectRepo,
		sessionRepo:         sessionRepo,
		retryConfigRepo:     retryConfigRepo,
		routingStrategyRepo: routingStrategyRepo,
		proxyRequestRepo:    proxyRequestRepo,
	}
}

// ServeHTTP routes admin requests
func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/admin")
	path = strings.TrimSuffix(path, "/")

	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	resource := parts[1]
	var id uint64
	if len(parts) > 2 && parts[2] != "" {
		id, _ = strconv.ParseUint(parts[2], 10, 64)
	}

	switch resource {
	case "providers":
		h.handleProviders(w, r, id)
	case "routes":
		h.handleRoutes(w, r, id)
	case "projects":
		h.handleProjects(w, r, id)
	case "sessions":
		h.handleSessions(w, r)
	case "retry-configs":
		h.handleRetryConfigs(w, r, id)
	case "routing-strategies":
		h.handleRoutingStrategies(w, r, id)
	case "requests":
		h.handleProxyRequests(w, r, id)
	default:
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

// Provider handlers
func (h *AdminHandler) handleProviders(w http.ResponseWriter, r *http.Request, id uint64) {
	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			provider, err := h.providerRepo.GetByID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "provider not found"})
				return
			}
			writeJSON(w, http.StatusOK, provider)
		} else {
			providers, err := h.providerRepo.List()
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, providers)
		}
	case http.MethodPost:
		var provider domain.Provider
		if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := h.providerRepo.Create(&provider); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, provider)
	case http.MethodPut:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		var provider domain.Provider
		if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		provider.ID = id
		if err := h.providerRepo.Update(&provider); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, provider)
	case http.MethodDelete:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		if err := h.providerRepo.Delete(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusNoContent, nil)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// Route handlers
func (h *AdminHandler) handleRoutes(w http.ResponseWriter, r *http.Request, id uint64) {
	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			route, err := h.routeRepo.GetByID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "route not found"})
				return
			}
			writeJSON(w, http.StatusOK, route)
		} else {
			routes, err := h.routeRepo.List()
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, routes)
		}
	case http.MethodPost:
		var route domain.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := h.routeRepo.Create(&route); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, route)
	case http.MethodPut:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		var route domain.Route
		if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		route.ID = id
		if err := h.routeRepo.Update(&route); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, route)
	case http.MethodDelete:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		if err := h.routeRepo.Delete(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusNoContent, nil)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// Project handlers
func (h *AdminHandler) handleProjects(w http.ResponseWriter, r *http.Request, id uint64) {
	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			project, err := h.projectRepo.GetByID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "project not found"})
				return
			}
			writeJSON(w, http.StatusOK, project)
		} else {
			projects, err := h.projectRepo.List()
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, projects)
		}
	case http.MethodPost:
		var project domain.Project
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := h.projectRepo.Create(&project); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, project)
	case http.MethodPut:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		var project domain.Project
		if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		project.ID = id
		if err := h.projectRepo.Update(&project); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, project)
	case http.MethodDelete:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		if err := h.projectRepo.Delete(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusNoContent, nil)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// Session handlers
func (h *AdminHandler) handleSessions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		sessions, err := h.sessionRepo.List()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, sessions)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// RetryConfig handlers
func (h *AdminHandler) handleRetryConfigs(w http.ResponseWriter, r *http.Request, id uint64) {
	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			config, err := h.retryConfigRepo.GetByID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "retry config not found"})
				return
			}
			writeJSON(w, http.StatusOK, config)
		} else {
			configs, err := h.retryConfigRepo.List()
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, configs)
		}
	case http.MethodPost:
		var config domain.RetryConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := h.retryConfigRepo.Create(&config); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, config)
	case http.MethodPut:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		var config domain.RetryConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		config.ID = id
		if err := h.retryConfigRepo.Update(&config); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, config)
	case http.MethodDelete:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		if err := h.retryConfigRepo.Delete(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusNoContent, nil)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// RoutingStrategy handlers
func (h *AdminHandler) handleRoutingStrategies(w http.ResponseWriter, r *http.Request, id uint64) {
	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			strategy, err := h.routingStrategyRepo.GetByProjectID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "routing strategy not found"})
				return
			}
			writeJSON(w, http.StatusOK, strategy)
		} else {
			strategies, err := h.routingStrategyRepo.List()
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, strategies)
		}
	case http.MethodPost:
		var strategy domain.RoutingStrategy
		if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := h.routingStrategyRepo.Create(&strategy); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, strategy)
	case http.MethodPut:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		var strategy domain.RoutingStrategy
		if err := json.NewDecoder(r.Body).Decode(&strategy); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		strategy.ID = id
		if err := h.routingStrategyRepo.Update(&strategy); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, strategy)
	case http.MethodDelete:
		if id == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
			return
		}
		if err := h.routingStrategyRepo.Delete(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusNoContent, nil)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

// ProxyRequest handlers
func (h *AdminHandler) handleProxyRequests(w http.ResponseWriter, r *http.Request, id uint64) {
	switch r.Method {
	case http.MethodGet:
		if id > 0 {
			req, err := h.proxyRequestRepo.GetByID(id)
			if err != nil {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "proxy request not found"})
				return
			}
			writeJSON(w, http.StatusOK, req)
		} else {
			// Get limit and offset from query params
			limit := 100
			offset := 0
			if l := r.URL.Query().Get("limit"); l != "" {
				limit, _ = strconv.Atoi(l)
			}
			if o := r.URL.Query().Get("offset"); o != "" {
				offset, _ = strconv.Atoi(o)
			}
			requests, err := h.proxyRequestRepo.List(limit, offset)
			if err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, requests)
		}
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
