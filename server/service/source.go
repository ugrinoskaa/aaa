package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/amukoski/aaa/model"
	"github.com/amukoski/aaa/service/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultSchema = "public"
	tmpUploadDir  = "uploads"
)

type SourceService struct {
	db *pgxpool.Pool
}

func NewSourceService(db *pgxpool.Pool) *SourceService {
	return &SourceService{db: db}
}

func (s *SourceService) All(ctx context.Context) ([]model.Source, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, type FROM sources`)
	if err != nil {
		return nil, errors.New("failed to retrieve sources")
	}
	defer rows.Close()

	sources := make([]model.Source, 0)
	for rows.Next() {
		var source model.Source
		if err = rows.Scan(&source.ID, &source.Name, &source.Type); err != nil {
			return nil, errors.New("failed to scan source row")
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (s *SourceService) Get(ctx context.Context, id int) (model.Source, []model.DatasetConfig, error) {
	source := model.Source{ID: id}
	err := s.db.QueryRow(ctx, `SELECT id, name, type, config FROM sources WHERE id = $1`, id).
		Scan(&source.ID, &source.Name, &source.Type, &source.Config)

	return source, source.Config.Datasets, err
}

type CreateSourceReq struct {
	Name     string
	Type     string
	Resource string
}

func (s *SourceService) Create(ctx context.Context, req CreateSourceReq) (int, error) {
	insertQuery := `
		INSERT INTO sources (name, type, config)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	if req.Type == string(model.POSTGRES) {
		config, id := model.SourceConfig{DatabaseURI: req.Resource}, 0
		err := s.db.QueryRow(ctx, insertQuery, req.Name, req.Type, config).Scan(&id)
		if err != nil {
			return 0, err
		}

		datasets, err := s.DiscoverDB(ctx, req.Resource)
		if err != nil {
			return 0, err
		}

		dsconfigs := make([]model.DatasetConfig, len(datasets))
		for idx, ds := range datasets {
			dsconfigs[idx] = ds.Config
		}

		config.Datasets = dsconfigs
		_, err = s.db.Exec(ctx, `UPDATE sources SET config = $1 WHERE id = $2;`, config, id)

		return id, err
	}

	if req.Type == string(model.CSV) {
		tx, err := s.db.Begin(ctx)
		if err != nil {
			return 0, err
		}

		datasets, err := s.DiscoverCSV(ctx, req.Resource)
		if err != nil {
			return 0, err
		}

		config, id, dsconfigs := model.SourceConfig{}, 0, make([]model.DatasetConfig, 0, len(datasets))
		err = s.db.QueryRow(ctx, insertQuery, req.Name, req.Type, config).Scan(&id)
		if err != nil {
			_ = tx.Rollback(ctx)
			return id, err
		}

		for _, ds := range datasets {
			table := fmt.Sprintf("sources_%d_%s", id, ds.Config.Table)
			columns := strings.ReplaceAll(strings.Join(ds.Config.Columns, ","), utils.ColumnSeparator, " ")
			createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %q (%s);", table, columns)

			if _, err = s.db.Exec(ctx, createTableQuery); err != nil {
				_ = tx.Rollback(ctx)
				return 0, err
			}

			filePath := filepath.Join(os.TempDir(), tmpUploadDir, req.Resource, ds.Config.Table)
			copySQL := fmt.Sprintf("COPY %q (%s) FROM STDIN WITH (FORMAT csv, HEADER true);",
				table, strings.Join(utils.ColumnNames(ds.Config.Columns), ","))

			file, err := os.Open(filePath)
			if err != nil {
				_ = tx.Rollback(ctx)
				return 0, err
			}

			_, err = tx.Conn().PgConn().CopyFrom(ctx, file, copySQL)
			_ = file.Close()

			if err != nil {
				_ = tx.Rollback(ctx)
				return 0, fmt.Errorf("failed to import csv: %v", err)
			}

			dsconfigs = append(dsconfigs, model.DatasetConfig{
				Schema:  defaultSchema,
				Table:   table,
				Columns: ds.Config.Columns,
			})
		}

		config.Datasets = dsconfigs
		if _, err = s.db.Exec(ctx, `UPDATE sources SET config = $1 WHERE id = $2;`, config, id); err != nil {
			_ = tx.Rollback(ctx)
			return 0, err
		}

		if err = tx.Commit(ctx); err != nil {
			return 0, err
		}

		return id, nil
	}

	return 0, errors.New("unsupported source type")
}

func (s *SourceService) Delete(ctx context.Context, id int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	source, datasets, err := s.Get(ctx, id)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if _, err = s.db.Exec(ctx, `DELETE FROM sources WHERE id = $1;`, id); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	if source.Type == model.CSV {
		tables := make([]string, len(datasets))
		for idx, ds := range datasets {
			tables[idx] = ds.Table
		}

		dropQuery := fmt.Sprintf("DROP TABLE IF EXISTS %s", strings.Join(tables, ","))
		if _, err = s.db.Exec(ctx, dropQuery); err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	return nil
}

func (s *SourceService) DiscoverDB(ctx context.Context, uri string) ([]model.Dataset, error) {
	datasets := make([]model.Dataset, 0)
	schemas := map[string][]model.DatasetConfig{}

	conn, err := pgxpool.Connect(ctx, uri)
	if err != nil {
		return datasets, errors.New("failed to connect to the database")
	}
	defer conn.Close()

	query := `
		SELECT table_schema, table_name, column_name, data_type
		FROM information_schema.columns
		WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		ORDER BY table_schema, table_name, ordinal_position;
	`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return datasets, errors.New("failed to query schemas, tables, and columns")
	}
	defer rows.Close()

	for rows.Next() {
		var schema, table, columnName, dataType string
		if err := rows.Scan(&schema, &table, &columnName, &dataType); err != nil {
			return datasets, errors.New("failed to scan row")
		}

		found := false
		for idx, info := range schemas[schema] {
			if info.Table == table {
				schemas[schema][idx].Columns = append(
					schemas[schema][idx].Columns,
					utils.FormatColumn(columnName, dataType),
				)
				found = true
				break
			}
		}

		if !found {
			schemas[schema] = append(schemas[schema],
				model.DatasetConfig{
					Schema:  schema,
					Table:   table,
					Columns: []string{utils.FormatColumn(columnName, dataType)},
				})
		}
	}

	for configs := range maps.Values(schemas) {
		for _, cfg := range configs {
			datasets = append(datasets, model.Dataset{Config: cfg})
		}
	}

	return datasets, nil
}
