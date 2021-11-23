package scheduler

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"micropairs/pkg/client/cryptocompare"
	"net/http"
	"time"
)

const (
	configFsyms = "scheduler.fsyms"
	configTsyms = "scheduler.tsyms"
	configSleep = "scheduler.sleep"
)

type Engine struct {
	cryptoClient *cryptocompare.Client
	log          *log.Logger
	ch           chan json.RawMessage
	tsyms, fsyms []string
	sleep        time.Duration
}

func NewScheduler(client *http.Client, logger *log.Logger) *Engine {
	return &Engine{
		log:          logger,
		cryptoClient: cryptocompare.NewClient(client, logger),
		ch:           make(chan json.RawMessage, 1),
		tsyms:        viper.GetStringSlice(configTsyms),
		fsyms:        viper.GetStringSlice(configFsyms),
		sleep:        viper.GetDuration(configSleep),
	}
}
func (e *Engine) Start(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			raw, err := e.cryptoClient.Request(ctx, e.fsyms, e.tsyms)
			if err != nil {
				e.log.WithError(err).Error("Scheduler problem")
			}
			if raw != nil {
				e.ch <- raw
			}
			time.Sleep(e.sleep)
		}
	}(ctx)
}
