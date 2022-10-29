package app

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/config"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/handlers"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/repositories"
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

	logger *log.Logger
	db     *sql.DB

	server  *gin.Engine
	handler *handlers.Handler
	gateway *gateways.Gateway
	repo    *repositories.Repo
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
	app.initDB()

	app.initRepo()
	app.initGateway()
	app.initHandler()

	app.initServer()
}

func (app *App) initServer() {
	app.server = gin.Default()

	app.server.GET("/health", app.handler.Health)
	app.server.GET("/photo", app.handler.GetRandPhoto)
	app.server.POST("/photo", app.handler.PostPhoto)

	err := app.server.Run(app.conf.Port)
	if err != nil {
		app.logger.Fatalf("error while init server: %v \n", err)
	}

}

func (app *App) initLogger() {
	app.logger = log.Default()
}

func (app *App) initHandler() {
	app.handler = handlers.NewHandler(app.gateway, app.repo, app.logger)
}

func (app *App) initGateway() {
	app.gateway = gateways.NewGateway(http.Client{}, app.conf.Token.UplashToken, app.logger)
}

func (app *App) initRepo() {
	app.repo = repositories.NewRepo(app.db, app.logger)
}

func (app *App) initDB() {
	db, err := sql.Open("sqlite3", "/home/tynrol/Code/GolandProjects/ITMO_IntelligentDataAnalysis/accessor-service/init/database.db")
	if err != nil {
		panic(err)
	}

	app.db = db
}
