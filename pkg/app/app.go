package app

import (
	"apous-films-rest-api/pkg/routes"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func (a *App) Initialize() {
	a.Router = gin.Default()
	routes.SetRoutes(a.Router)
}

func (a *App) Run() {
	a.Router.Run()
}
