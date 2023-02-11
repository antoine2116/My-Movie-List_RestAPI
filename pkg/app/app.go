package app

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/pkg/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
}

func (a *App) Initialize(c config.Configuration) {
	a.Router = gin.Default()

	routes.SetRoutes(a.Router)
}

func (a *App) Run(port int) {
	a.Router.Run(fmt.Sprintf(":%d", port))
}
