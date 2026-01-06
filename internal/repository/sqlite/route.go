package sqlite

import (
	"database/sql"
	"time"

	"github.com/Bowl42/maxx-next/internal/domain"
)

type RouteRepository struct {
	db *DB
}

func NewRouteRepository(db *DB) *RouteRepository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) Create(route *domain.Route) error {
	now := time.Now()
	route.CreatedAt = now
	route.UpdatedAt = now

	isEnabled := 0
	if route.IsEnabled {
		isEnabled = 1
	}

	result, err := r.db.db.Exec(
		`INSERT INTO routes (created_at, updated_at, is_enabled, project_id, client_type, provider_id, position, retry_config_id, model_mapping) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		route.CreatedAt, route.UpdatedAt, isEnabled, route.ProjectID, route.ClientType, route.ProviderID, route.Position, route.RetryConfigID, toJSON(route.ModelMapping),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	route.ID = uint64(id)
	return nil
}

func (r *RouteRepository) Update(route *domain.Route) error {
	route.UpdatedAt = time.Now()
	isEnabled := 0
	if route.IsEnabled {
		isEnabled = 1
	}
	_, err := r.db.db.Exec(
		`UPDATE routes SET updated_at = ?, is_enabled = ?, project_id = ?, client_type = ?, provider_id = ?, position = ?, retry_config_id = ?, model_mapping = ? WHERE id = ?`,
		route.UpdatedAt, isEnabled, route.ProjectID, route.ClientType, route.ProviderID, route.Position, route.RetryConfigID, toJSON(route.ModelMapping), route.ID,
	)
	return err
}

func (r *RouteRepository) Delete(id uint64) error {
	_, err := r.db.db.Exec(`DELETE FROM routes WHERE id = ?`, id)
	return err
}

func (r *RouteRepository) GetByID(id uint64) (*domain.Route, error) {
	row := r.db.db.QueryRow(`SELECT id, created_at, updated_at, is_enabled, project_id, client_type, provider_id, position, retry_config_id, model_mapping FROM routes WHERE id = ?`, id)
	return r.scanRoute(row)
}

func (r *RouteRepository) List() ([]*domain.Route, error) {
	rows, err := r.db.db.Query(`SELECT id, created_at, updated_at, is_enabled, project_id, client_type, provider_id, position, retry_config_id, model_mapping FROM routes ORDER BY position`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*domain.Route
	for rows.Next() {
		route, err := r.scanRouteRows(rows)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, rows.Err()
}

func (r *RouteRepository) scanRoute(row *sql.Row) (*domain.Route, error) {
	var route domain.Route
	var isEnabled int
	var mappingJSON string
	err := row.Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt, &isEnabled, &route.ProjectID, &route.ClientType, &route.ProviderID, &route.Position, &route.RetryConfigID, &mappingJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	route.IsEnabled = isEnabled == 1
	route.ModelMapping = fromJSON[map[string]string](mappingJSON)
	return &route, nil
}

func (r *RouteRepository) scanRouteRows(rows *sql.Rows) (*domain.Route, error) {
	var route domain.Route
	var isEnabled int
	var mappingJSON string
	err := rows.Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt, &isEnabled, &route.ProjectID, &route.ClientType, &route.ProviderID, &route.Position, &route.RetryConfigID, &mappingJSON)
	if err != nil {
		return nil, err
	}
	route.IsEnabled = isEnabled == 1
	route.ModelMapping = fromJSON[map[string]string](mappingJSON)
	return &route, nil
}
