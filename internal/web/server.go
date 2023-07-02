package web

import (
	"fmt"

	"github.com/amleonc/tabula/config"
	"github.com/amleonc/tabula/internal/web/internal/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Router struct {
	server *echo.Echo
	port   int
}

func (r *Router) Start() {
	mountPublicRoutes(r)
	r.server.Logger.Fatal(r.server.Start(fmt.Sprintf(":%d", r.port)))
}

func New() Router {
	e := echo.New()
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 2 << 10,
		LogLevel:  log.ERROR,
	}))
	p := config.ServerPort()
	return Router{
		server: e,
		port:   p,
	}
}

func mountPublicRoutes(r *Router) {
	r.server.POST("/signup", controllers.Signup)
	r.server.POST("/login", controllers.Login)
}
