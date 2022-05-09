package src

import (
	"github.com/steevepypo/todoback/settings"
	"github.com/steevepypo/todoback/src/controllers"
)

func Run() {
	server := controllers.Server{}
	settings.LoadEnv()
	server.InitDB()
	server.InitRouter()
}
