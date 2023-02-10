package app

import "github.com/gin-gonic/gin"

type App struct {
	Router *gin.Engine
}

func (a *App) Initialize() {
	a.Router = gin.Default()

}
