package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"micropairs/internal/scheduler"
	"micropairs/pkg/client/cryptocompare"
	"micropairs/pkg/client/db"
	"micropairs/pkg/helper"
)

const (
	configPort = "server.port"
)
const (
	queryFsyms = "fsyms"
	queryTsyms = "tsyms"
)

type Server struct {
	r         *gin.Engine
	log       *log.Logger
	client    *cryptocompare.Client
	scheduler *scheduler.Engine
	db        *db.Client
}

func NewCryptoServer(client *http.Client) *Server {
	logger := log.New()
	return &Server{
		r:      gin.Default(),
		log:    logger,
		client: cryptocompare.NewClient(client, logger),
		db:     db.NewClient(),
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
	if err := s.db.Connect(); err != nil {
		return err
	}
	s.db.CreateSchema()

	s.r.GET("/price", s.handleRequest)
	err := s.r.Run(viper.GetString(configPort))
	go s.UpdateOnScheduler()
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
		s.log.WithError(err).Error("err on client request")
		s.log.Info("trying to get last data from db")
		if raw, err = s.GetDataFromDB(fsyms, tsyms); err != nil {
			c.Status(http.StatusBadGateway)
			return
		}
	}
	c.JSON(http.StatusOK, raw)
}

func (s *Server) GetDataFromDB(fsyms, tsyms []string) (json.RawMessage, error) {
	raw, err := s.db.GetLastRawJSON()
	if err != nil {
		s.log.WithError(err).Error("Got an error from db")
		return nil, err
	}
	res, err := helper.GetPairsFromRawJSON(fsyms, tsyms, raw)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Server) UpdateOnScheduler() {
	for json := range s.scheduler.Ch {
		if err := s.db.InsertRawJSON(json); err != nil {
			continue
		}
	}
}
