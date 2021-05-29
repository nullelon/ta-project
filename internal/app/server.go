package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
	"ta-project-go/assets"
)

type Server struct {
	r *gin.Engine

	MarkerService *MarkerService

	addr string
}

func NewServer(addr string) (*Server, error) {
	server := &Server{
		r:             gin.Default(),
		MarkerService: NewMarkerService(),
		addr:          addr,
	}
	server.configureRouter()
	return server, nil
}

func (s Server) Start() error {
	err := s.r.Run(s.addr)
	if err != nil {
		return err
	}
	return nil
}

func (s Server) configureRouter() {
	s.r.Use(cors.Default())
	s.r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", assets.IndexHtml)
	})
	s.r.GET("/api/info", func(c *gin.Context) {
		symbol, ok := c.GetQuery("symbol")
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"error": "symbol is required parameter"})
			return
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "300"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "limit must be a number"})
			return
		}
		timeframe := c.DefaultQuery("timeframe", "5m")

		candlesticks, err := s.MarkerService.Get(symbol, timeframe, limit)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
			return
		}
		c.JSON(200, candlesticks)
	})
}
