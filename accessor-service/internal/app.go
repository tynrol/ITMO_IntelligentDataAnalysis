package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/config"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/gateways"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/handlers"
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

	logger  *log.Logger
	server  *gin.Engine
	handler *handlers.Handler
	gateway *gateways.Gateway
	//service *services.Service
	//repo    *repo.Repo
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

	app.initGateway()
	app.initHandler()
	app.initService()
	app.initRepo()

	app.initServer()
}

func (app *App) initServer() {
	app.server = gin.Default()

	app.server.GET("/health", app.handler.Health)
	app.server.GET("/photo", app.handler.GetRandPhoto)

	err := app.server.Run(app.conf.Port)
	if err != nil {
		app.logger.Fatalf("error while init server: %v \n", err)
	}

}

func (app *App) initLogger() {
	app.logger = log.Default()
}

func (app *App) initHandler() {
	app.handler = handlers.NewHandler(app.gateway, app.logger)
}

func (app *App) initService() {

}

func (app *App) initGateway() {
	app.gateway = gateways.NewGateway(http.Client{}, app.conf.Token.UplashToken, app.logger)
}

func (app *App) initRepo() {

}

func (app *App) initDB() {

}
