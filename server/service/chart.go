package service

import (
	"context"
	"fmt"
	"github.com/amukoski/aaa/model"
	"github.com/amukoski/aaa/service/utils"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Chart interface {
	Schema() model.ChartSchema
	Render(name string, groups [][]string, values [][]float64) (any, error)
}

type ChartService struct {
	db       *pgxpool.Pool
	sources  *SourceService
	datasets *DatasetService
	registry map[model.ChartType]Chart
}

func NewChartService(db *pgxpool.Pool, src *SourceService, ds *DatasetService, charts ...Chart) *ChartService {
	registry := make(map[model.ChartType]Chart)
	for _, chart := range charts {
		schema := chart.Schema()
		registry[schema.Type] = chart
	}

	return &ChartService{
		db:       db,
		sources:  src,
		datasets: ds,
		registry: registry,
	}
}

func (s *ChartService) All(ctx context.Context) ([]model.Chart, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, dataset_id, type FROM charts`)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve charts: %w", err)
	}
	defer rows.Close()

	charts := make([]model.Chart, 0)
	for rows.Next() {
		var chart model.Chart
		if err = rows.Scan(&chart.ID, &chart.Name, &chart.DatasetID, &chart.Type); err != nil {
			return nil, fmt.Errorf("failed to scan chart row: %w", err)
		}
		charts = append(charts, chart)
	}

	return charts, nil
}

func (s *ChartService) AllTypes() []model.ChartType {
	return model.SupportedChartTypes
}

func (s *ChartService) Get(ctx context.Context, id int) (model.Chart, error) {
	query := `SELECT id, name, dataset_id, type, config FROM charts WHERE id = $1`

	var chart model.Chart
	err := s.db.QueryRow(ctx, query, id).
		Scan(&chart.ID, &chart.Name, &chart.DatasetID, &chart.Type, &chart.Config)
	if err != nil {
		return chart, fmt.Errorf("failed to retrieve chart: %w", err)
	}

	return chart, nil
}

func (s *ChartService) GetType(ctype string) (model.ChartSchema, error) {
	chart, found := s.registry[model.ChartType(ctype)]
	if !found {
		return model.ChartSchema{}, fmt.Errorf("unknown chart type: %s", ctype)
	}

	return chart.Schema(), nil
}

type CreateChartReq struct {
	DatasetID  int
	Name       string
	Type       string
	Dimensions []string
	Metrics    []string
	Filters    []string
}

func (s *ChartService) Create(ctx context.Context, req CreateChartReq) (int, error) {
	query := `
		INSERT INTO charts (dataset_id, name, type, config)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	config := model.ChartConfig{
		Dimensions: req.Dimensions,
		Metrics:    req.Metrics,
		Filters:    req.Filters,
	}

	var id int
	err := s.db.QueryRow(ctx, query, req.DatasetID, req.Name, req.Type, config).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create chart: %w", err)
	}

	return id, nil
}

func (s *ChartService) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM charts WHERE id = $1;`
	_, err := s.db.Exec(ctx, query, id)
	return err
}

type ValidateChartReq struct {
	DatasetID  int
	Name       string
	Type       string
	Dimensions []string
	Metrics    []string
	Filters    []string
}

func (s *ChartService) Validate(ctx context.Context, req ValidateChartReq) (any, error) {
	chart, found := s.registry[model.ChartType(req.Type)]
	if !found {
		return model.ChartSchema{}, fmt.Errorf("unknown chart type: %s", req.Type)
	}

	var result map[string]interface{}

	dataset, err := s.datasets.Get(ctx, req.DatasetID)
	if err != nil {
		return result, fmt.Errorf("failed to retrieve dataset: %w", err)
	}

	source, _, err := s.sources.Get(ctx, dataset.SourceID)
	if err != nil {
		return result, fmt.Errorf("failed to retrieve source: %w", err)
	}

	conn := s.db
	if source.Type == model.POSTGRES {
		conn, err = pgxpool.Connect(ctx, source.Config.DatabaseURI)
		if err != nil {
			return result, fmt.Errorf("failed to connect to database: %w", err)
		}
		defer conn.Close()
	}

	table := fmt.Sprintf("%s.%s", dataset.Config.Schema, dataset.Config.Table)
	query := utils.BuildSQLQuery(table, req.Dimensions, req.Metrics, req.Filters, dataset.Config.Columns)

	fmt.Println(fmt.Sprintf("%v", query))

	groups, values, err := s.perform(ctx, conn, query, len(req.Dimensions), len(req.Metrics))
	if err != nil {
		return result, err
	}

	return chart.Render(req.Name, groups, values)
}

func (s *ChartService) perform(ctx context.Context, conn *pgxpool.Pool, query string, dimensions int, metrics int) ([][]string, [][]float64, error) {
	groups := make([][]string, dimensions)
	values := make([][]float64, metrics)

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate chart: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		scans, err := rows.Values()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan row: %w", err)
		}

		for idx, scan := range scans {
			if idx < dimensions {
				if scan == nil {
					scan = "N/A"
				}

				if val, err := utils.ToFloat64(scan); err == nil {
					scan = val
				}

				groups[idx] = append(groups[idx], fmt.Sprintf("%v", scan))
				continue
			}

			if idx-dimensions < metrics {
				fvalue, _ := utils.ToFloat64(scan)
				values[idx-dimensions] = append(values[idx-dimensions], fvalue)
			}
		}
	}

	return groups, values, nil
}

func (s *ChartService) Run(ctx context.Context, id int) (any, error) {
	var result any
	var chart model.Chart

	query := `SELECT name, type, dataset_id, config FROM charts WHERE id = $1`
	err := s.db.QueryRow(ctx, query, id).Scan(&chart.Name, &chart.Type, &chart.DatasetID, &chart.Config)
	if err != nil {
		return result, fmt.Errorf("failed to retrieve chart: %w", err)
	}

	if result, err = s.Validate(ctx, ValidateChartReq{
		DatasetID:  chart.DatasetID,
		Name:       chart.Name,
		Type:       string(chart.Type),
		Dimensions: chart.Config.Dimensions,
		Metrics:    chart.Config.Metrics,
		Filters:    chart.Config.Filters,
	}); err != nil {
		return result, fmt.Errorf("failed to validate chart: %w", err)
	}

	return result, nil
}
