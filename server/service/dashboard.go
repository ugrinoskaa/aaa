package service

import (
	"context"
	"errors"

	"github.com/amukoski/aaa/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DashboardService struct {
	db *pgxpool.Pool
}

func NewDashboardService(db *pgxpool.Pool) *DashboardService {
	return &DashboardService{db: db}
}

func (s *DashboardService) All(ctx context.Context) ([]model.Dashboard, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, grid FROM dashboards`)
	if err != nil {
		return nil, errors.New("failed to retrieve dashboards")
	}
	defer rows.Close()

	dashboards := make([]model.Dashboard, 0)
	for rows.Next() {
		var dashboard model.Dashboard
		if err = rows.Scan(&dashboard.ID, &dashboard.Name, &dashboard.Grid); err != nil {
			return nil, errors.New("failed to scan dashboard row")
		}
		dashboards = append(dashboards, dashboard)
	}

	return dashboards, nil
}

func (s *DashboardService) Get(ctx context.Context, id int) (model.Dashboard, error) {
	query := `SELECT id, name, grid FROM dashboards WHERE id = $1`

	var dashboard model.Dashboard
	err := s.db.QueryRow(ctx, query, id).Scan(&dashboard.ID, &dashboard.Name, &dashboard.Grid)
	if err != nil {
		return dashboard, errors.New("dashboard not found")
	}

	return dashboard, nil
}

type CreateDashboardReq struct {
	Name string
	Grid []map[string]any
}

func (s *DashboardService) Create(ctx context.Context, req CreateDashboardReq) (int, error) {
	query := `INSERT INTO dashboards (name, grid) VALUES ($1, $2) RETURNING id;`

	var id int
	err := s.db.QueryRow(ctx, query, req.Name, req.Grid).Scan(&id)
	if err != nil {
		return 0, errors.New("failed to insert dashboard")
	}

	return id, nil
}

func (s *DashboardService) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM dashboards WHERE id = $1;`
	_, err := s.db.Exec(ctx, query, id)
	return err
}
