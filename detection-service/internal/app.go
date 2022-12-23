package app

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/detection-service/config"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/detection-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/detection-service/internal/handlers"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
	name   string

	conf *config.Config

	logger     *log.Logger
	db         *sql.DB
	metricsMux *http.ServeMux

	server  *gin.Engine
	handler *handlers.Handler
	gateway *gateways.Gateway
}

func NewApp(conf *config.Config) *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:    ctx,
		cancel: cancel,
		name:   conf.App.Name,
		conf:   conf,
	}
}

func (app *App) Start() {
	app.init()

	app.logger.Print("service is started")
}

func (app *App) Close() {
	app.cancel()

	app.logger.Print("service is stopped")
}

func (app *App) init() {
	rand.Seed(time.Now().UnixNano())
	app.initLogger()
	app.initMetrics()
	app.initModel()

	app.initRepo()
	app.initGateway()
	app.initHandler()

	app.initServer()
}

func (app *App) initServer() {
	app.server = gin.Default()

	app.server.Use(CORSMiddleware())

	app.server.GET("/metrics", prometheusHandler())
	app.server.GET("/health", app.handler.Health)
	app.server.POST("/detect", app.handler.Detect)

	err := app.server.Run(app.conf.HTTPConfig.Port)
	if err != nil {
		app.logger.Fatalf("error while init server: %v \n", err)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (app *App) initMetrics() {
	usersRegistered := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "users_registered",
		})
	prometheus.MustRegister(usersRegistered)

	usersOnline := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_online",
		})
	prometheus.MustRegister(usersOnline)

	requestProcessingTimeSummaryMs := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "request_processing_time_summary_ms",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(requestProcessingTimeSummaryMs)

	requestProcessingTimeHistogramMs := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "request_processing_time_histogram_ms",
			Buckets: prometheus.LinearBuckets(0, 10, 20),
		})
	prometheus.MustRegister(requestProcessingTimeHistogramMs)

	go func() {
		for {
			usersRegistered.Inc() // or: Add(5)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		for {
			for i := 0; i < 10000; i++ {
				usersOnline.Set(float64(i)) // or: Inc(), Dec(), Add(5), Dec(5)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()

	go func() {
		src := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(src)
		for {
			obs := float64(100 + rnd.Intn(30))
			requestProcessingTimeSummaryMs.Observe(obs)
			requestProcessingTimeHistogramMs.Observe(obs)
			time.Sleep(10 * time.Millisecond)
		}
	}()

}

func (app *App) initLogger() {
	app.logger = log.Default()
}

func (app *App) initModel() {
	//flag.StringVar(&app.model, "model", "/home/tynrol/Code/GoLandProjects/ITMO_IntelligentDataAnalysis/detection-service/models/model.pt", "Model weights for inference")
}

func (app *App) initGateway() {
	app.gateway = gateways.NewGateway(http.Client{}, app.logger)
}

func (app *App) initRepo() {
}

func (app *App) initHandler() {
	app.handler = handlers.NewHandler(app.gateway, app.logger)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
