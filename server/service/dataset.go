package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/amukoski/aaa/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DatasetService struct {
	db      *pgxpool.Pool
	sources *SourceService
}

func NewDatasetService(db *pgxpool.Pool, src *SourceService) *DatasetService {
	return &DatasetService{
		db:      db,
		sources: src,
	}
}

func (s *DatasetService) All(ctx context.Context) ([]model.Dataset, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, source_id, config FROM datasets`)
	if err != nil {
		return nil, errors.New("failed to retrieve datasets")
	}
	defer rows.Close()

	datasets := make([]model.Dataset, 0)
	for rows.Next() {
		var dataset model.Dataset
		if err = rows.Scan(&dataset.ID, &dataset.Name, &dataset.SourceID, &dataset.Config); err != nil {
			return nil, errors.New("failed to scan dataset row")
		}
		datasets = append(datasets, dataset)
	}

	return datasets, nil
}

func (s *DatasetService) Get(ctx context.Context, id int) (model.Dataset, error) {
	query := `SELECT id, name, source_id, config FROM datasets WHERE id = $1`

	var dataset model.Dataset
	err := s.db.QueryRow(ctx, query, id).Scan(&dataset.ID, &dataset.Name, &dataset.SourceID, &dataset.Config)
	if err != nil {
		return dataset, errors.New("dataset not found")
	}

	return dataset, nil
}

type CreateDatasetReq struct {
	Name           string
	SourceID       int
	DatabaseSchema string
	DatabaseTable  string
}

func (s *DatasetService) Create(ctx context.Context, req CreateDatasetReq) (int, error) {
	_, datasets, err := s.sources.Get(ctx, req.SourceID)
	if err != nil {
		return 0, fmt.Errorf("failed to get source %s: %w", req.SourceID, err)
	}

	config := model.DatasetConfig{
		Schema:  req.DatabaseSchema,
		Table:   req.DatabaseTable,
		Columns: []string{},
	}

	for _, ds := range datasets {
		if ds.Table == req.DatabaseTable && ds.Schema == req.DatabaseSchema {
			config.Columns = ds.Columns
			break
		}
	}

	query := `
		INSERT INTO datasets (name, source_id, config)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id int
	err = s.db.QueryRow(ctx, query, req.Name, req.SourceID, config).Scan(&id)
	if err != nil {
		return 0, errors.New("failed to insert dataset")
	}

	return id, nil
}

func (s *DatasetService) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM datasets WHERE id = $1;`
	_, err := s.db.Exec(ctx, query, id)
	return err
}
