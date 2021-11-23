package cryptocompare

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	configURI     = "crypto.URI"
	configTimeout = "crypto.timeout"
)
const (
	fparam = "fsyms"
	tparam = "tsyms"
)

type Client struct {
	client  *http.Client
	log     *log.Logger
	URL     string
	timeout time.Duration
}

func NewClient(c *http.Client, log *log.Logger) *Client {
	return &Client{
		client:  c,
		log:     log,
		URL:     viper.GetString(configURI),
		timeout: viper.GetDuration(configTimeout),
	}
}

func (s *Client) Request(ctx context.Context, fsyms, tsyms []string) (json.RawMessage, error) {
	ctxLocal, cancelCtx := context.WithTimeout(ctx, s.timeout)
	defer cancelCtx()

	req, err := http.NewRequestWithContext(ctxLocal, "GET", s.URL, nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Set(fparam, strings.Join(fsyms, ","))
	params.Set(tparam, strings.Join(tsyms, ","))
	req.URL.RawQuery = params.Encode()

	s.log.WithField("query:", req.URL).Info("builtQuery")

	r, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	s.log.Info("Successful request")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = r.Body.Close()
	}()
	raw := make(json.RawMessage, len(body))
	copy(raw, body)

	return raw, nil
}
