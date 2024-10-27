package service

import (
	"context"
	"encoding/csv"
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

func (s *SourceService) DiscoverCSV(ctx context.Context, uploadID string) ([]model.Dataset, error) {
	datasets := make([]model.Dataset, 0)
	uploadDir := filepath.Join(os.TempDir(), tmpUploadDir, uploadID)
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		return datasets, err
	}

	for _, entry := range entries {
		if ctx.Err() != nil {
			return datasets, ctx.Err()
		}

		if entry.IsDir() {
			continue
		}

		err = func() error {
			file, err := os.Open(filepath.Join(uploadDir, entry.Name()))
			if err != nil {
				return err
			}
			defer file.Close()

			columns, err := resolve(file)
			if err != nil {
				return err
			}

			datasets = append(datasets, model.Dataset{
				Config: model.DatasetConfig{
					Table:   entry.Name(),
					Schema:  uploadID,
					Columns: columns,
				},
			})

			return err
		}()
		if err != nil {
			return datasets, err
		}
	}

	return datasets, err
}

func (s *SourceService) Upload(ctx context.Context, files []*multipart.FileHeader) ([]model.Dataset, error) {
	uploadID := uuid.NewString()
	uploadDir := filepath.Join(os.TempDir(), tmpUploadDir, uploadID)
	datasets := make([]model.Dataset, 0, len(files))

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("unable to create upload directory: %v", err)
	}

	for _, file := range files {
		if ctx.Err() != nil {
			return datasets, ctx.Err()
		}

		filename := filepath.Clean(file.Filename)
		filename = strings.TrimSpace(strings.ToLower(filename))
		filename = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(filename, "_")

		err := func() error {
			incoming, err := file.Open()
			if err != nil {
				return err
			}
			defer incoming.Close()

			// save to: /tmp/uploads/{uploadID}/{filename}
			destPath := filepath.Join(uploadDir, filename)
			dest, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer dest.Close()

			if _, err := io.Copy(dest, incoming); err != nil {
				return err
			}

			readFile, err := os.Open(destPath)
			if err != nil {
				return err
			}
			defer readFile.Close()

			columns, err := resolve(readFile)
			if err != nil {
				return err
			}

			datasets = append(datasets, model.Dataset{
				Config: model.DatasetConfig{
					Table:   filename,
					Schema:  uploadID,
					Columns: columns,
				},
			})

			return err
		}()
		if err != nil {
			return datasets, err
		}
	}

	go func(dir string) {
		time.Sleep(10 * time.Minute)
		_ = os.RemoveAll(dir)
	}(uploadDir)

	return datasets, nil
}

func resolve(file io.Reader) ([]string, error) {
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	var records [][]string
	for i := 0; i < 100; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		records = append(records, row)
	}

	colValues := make([][]string, len(headers))
	for _, row := range records {
		for i := range headers {
			if i < len(row) {
				colValues[i] = append(colValues[i], row[i])
			}
		}
	}

	columns := make([]string, len(headers))
	for i, name := range headers {
		columns[i] = utils.FormatColumn(name, detectType(colValues[i]))
	}

	return columns, nil
}

func detectType(values []string) string {
	isNumber, isDate := true, true

	for _, val := range values {
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			isNumber = false
		}
		if _, err := time.Parse("2006-01-02", val); err != nil {
			isDate = false
		}
	}
	switch {
	case isNumber:
		return "numeric"
	case isDate:
		return "date"
	default:
		return "text"
	}
}
