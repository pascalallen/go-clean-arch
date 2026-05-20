package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pascalallen/go-clean-arch/internal/app/application/http/middleware"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
)

const v1 = "/api/v1"

type Router struct {
	engine *gin.Engine
}

func NewRouter() Router {
	return Router{
		engine: gin.Default(),
	}
}

// UseLogger installs the application logging middleware using the provided base logger.
func (r Router) UseLogger(l logger.Logger) {
	r.engine.Use(middleware.LoggerMiddleware(l))
}

func (r Router) Serve(port string) {
	if err := r.engine.Run(port); err != nil {
		log.Fatal(err)
	}
}
