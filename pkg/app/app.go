package app

import (
	"apous-films-rest-api/config"
	"apous-films-rest-api/pkg/routes"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router *gin.Engine
	DB     *mongo.Client
}

func (a *App) Initialize(c config.Configuration) {
	// Router
	a.Router = gin.Default()
	routes.SetRoutes(a.Router)

	// Database
	a.DB = config.GetDBClient(c.Database.URI)
}

func (a *App) Run(port int) {
	a.Router.Run(fmt.Sprintf(":%d", port))
}
