package scheduler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"micropairs/pkg/client/cryptocompare"
)

const (
	configFsyms = "scheduler.fsyms"
	configTsyms = "scheduler.tsyms"
	configSleep = "scheduler.sleep"
)

type Engine struct {
	cryptoClient *cryptocompare.Client
	log          *log.Logger
	Ch           chan json.RawMessage
	tsyms, fsyms []string
	sleep        time.Duration
}

func NewScheduler(client *http.Client, logger *log.Logger) *Engine {
	return &Engine{
		log:          logger,
		cryptoClient: cryptocompare.NewClient(client, logger),
		Ch:           make(chan json.RawMessage, 1),
		tsyms:        viper.GetStringSlice(configTsyms),
		fsyms:        viper.GetStringSlice(configFsyms),
		sleep:        viper.GetDuration(configSleep),
	}
}
func (e *Engine) Start(ctx context.Context) {
	ticker := time.NewTicker(e.sleep)
	go func(ctx context.Context, ticker *time.Ticker) {
		for range ticker.C {
			raw, err := e.cryptoClient.Request(ctx, e.fsyms, e.tsyms)
			if err != nil {
				e.log.WithError(err).Error("Scheduler problem")
			}
			if raw != nil {
				e.Ch <- raw
			}
		}
	}(ctx, ticker)
}
