package server

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"micropairs/internal/scheduler"
	"micropairs/pkg/client/cryptocompare"
	"net/http"
)

const (
	configPort = "server.port"
	queryFsyms = "fsyms"
	queryTsyms = "tsyms"
)

type Server struct {
	r         *gin.Engine
	log       *log.Logger
	client    *cryptocompare.Client
	scheduler *scheduler.Engine
}

func NewCryptoServer(client *http.Client) *Server {
	logger := log.New()
	return &Server{
		r:      gin.Default(),
		log:    logger,
		client: cryptocompare.NewClient(client, logger),
	}
}
func (s *Server) WithScheduler(client *http.Client) {
	s.scheduler = scheduler.NewScheduler(client, s.log)
}

// RunScheduler runs request with interval and writes data to db.
func (s *Server) RunScheduler() {
	if s.scheduler != nil {
		s.scheduler.Start(context.Background())
	}
}

func (s *Server) Run() error {
	s.r.Use(
		gin.Recovery(),
		gin.Logger(),
	)
	s.r.GET("/price", s.handleRequest)
	err := s.r.Run(viper.GetString(configPort))
	return err
}

func (s *Server) handleRequest(c *gin.Context) {
	fsyms, ok := c.GetQueryArray(queryFsyms)
	if !ok {
		s.log.Error("no fsyms provided")
		c.Status(http.StatusBadRequest)
		return
	}
	tsyms, ok := c.GetQueryArray(queryTsyms)
	if !ok {
		s.log.Error("no tsyms provided")
		c.Status(http.StatusBadRequest)
		return
	}
	raw, err := s.client.Request(c, fsyms, tsyms)
	if err != nil {
		s.log.Error("err on client request", err)
		c.Status(http.StatusBadGateway)
		return
	}
	c.JSON(http.StatusOK, raw)
}
