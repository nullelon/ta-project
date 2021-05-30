package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
)

type Server struct {
	r *gin.Engine

	MarkerService *MarkerService

	addr string
}

func NewServer(addr string) (*Server, error) {
	marketService := NewMarketService()
	err := marketService.Start()
	if err != nil {
		return nil, err
	}
	server := &Server{
		r:             gin.Default(),
		MarkerService: marketService,
		addr:          addr,
	}
	server.configureRouter()
	return server, nil
}

func (s *Server) Start() error {
	err := s.r.Run(s.addr)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) configureRouter() {
	s.r.Use(cors.Default())

	s.r.Static("/static", "assets/static")
	s.r.StaticFile("/", "assets/index.html")
	s.r.StaticFile("/graphs", "assets/graphs.html")

	s.r.GET("/api/info", func(c *gin.Context) {
		symbol, ok := c.GetQuery("symbol")
		if !ok {
			c.AbortWithStatusJSON(400, gin.H{"error": "symbol is required parameter"})
			return
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "300"))
		if err != nil || limit <= 0 || limit > 1000 {
			c.AbortWithStatusJSON(400, gin.H{"error": "limit must be a positive number < 1000"})
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
