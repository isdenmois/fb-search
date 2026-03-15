package views

import (
	"fb-search/views/controllers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HttpServer struct {
	pool  *pgxpool.Pool
	r     *gin.Engine
	ctrls *[]controllers.Controller
}

func NewHttpServer(pool *pgxpool.Pool, ctrls *[]controllers.Controller) *HttpServer {
	r := gin.Default()

	return &HttpServer{
		pool:  pool,
		r:     r,
		ctrls: ctrls,
	}
}

func (h *HttpServer) Run() error {
	for _, ctrl := range *h.ctrls {
		ctrl.Bind(h.r)
	}

	return h.r.Run()
}
