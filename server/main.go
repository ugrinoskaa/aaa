package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/amukoski/aaa/api"
	"github.com/amukoski/aaa/service"
	"github.com/amukoski/aaa/service/render"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	defaultDatabaseURL = "postgres://admin:admin@localhost:5432/aaa"
	defaultHttpPort    = "8080"
)

func main() {
	ctx := context.Background()
	logger := log.Default()

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = defaultDatabaseURL
		logger.Print("DATABASE_URL environment variable not set")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = defaultHttpPort
		logger.Print("HTTP_PORT environment variable not set")
	}

	db, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	barChart, barErr := render.NewBarChart()
	pieChart, pieErr := render.NewPieChart()
	lineChart, lineErr := render.NewLineChart()
	scatterChart, scatterErr := render.NewScatterChart()
	heatmapChart, heatmapErr := render.NewHeatmapChart()
	sankeyChart, sankeyErr := render.NewSankeyChart()

	if err = errors.Join(barErr, pieErr, lineErr, scatterErr, heatmapErr, sankeyErr); err != nil {
		logger.Fatal(err)
	}

	registry := []service.Chart{barChart, pieChart, lineChart, scatterChart, heatmapChart, sankeyChart}

	sources := service.NewSourceService(db)
	datasets := service.NewDatasetService(db, sources)
	charts := service.NewChartService(db, sources, datasets, registry...)
	dashboards := service.NewDashboardService(db)
	handler := api.Handler{
		Logger:    logger,
		Sources:   sources,
		Datasets:  datasets,
		Charts:    charts,
		Dashboard: dashboards,
	}

	app := fiber.New()
	handler.RegisterRoutes(app.Group("/api"))

	if err = app.Listen(fmt.Sprintf(":%v", httpPort)); err != nil {
		logger.Fatal(err)
	}
}
