package elasticsearch

import (
	"context"
	"encoding/base64"
	"encoding/json"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/drivers"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

type Connector struct {
	stack         string
	serviceConfig drivers.ServiceConfig
	config        Config
	client        *elastic.Client
	logger        logging.Logger
}

func (connector *Connector) Stop(_ context.Context) error {
	connector.client.Stop()
	return nil
}

func (connector *Connector) Start(ctx context.Context) error {
	options := []elastic.ClientOptionFunc{
		elastic.SetURL(connector.config.Endpoint),
	}
	if connector.config.Authentication != nil {
		options = append(options, elastic.SetBasicAuth(connector.config.Authentication.Username, connector.config.Authentication.Password))
	}
	if connector.serviceConfig.Debug {
		options = append(options,
			elastic.SetErrorLog(newLogger(connector.logger.Errorf)),
			elastic.SetInfoLog(newLogger(connector.logger.Infof)),
			elastic.SetTraceLog(newLogger(connector.logger.Debugf)),
		)
	}
	var err error
	connector.client, err = elastic.NewClient(options...)
	if err != nil {
		return errors.Wrap(err, "building es client")
	}

	return nil
}

func (connector *Connector) Client() *elastic.Client {
	return connector.client
}

func (connector *Connector) Accept(ctx context.Context, logs ...ingester.LogWithModule) ([]error, error) {

	bulk := connector.client.Bulk().Refresh("true")
	for _, log := range logs {
		doc := struct {
			ID      string          `json:"id"`
			Stack   string          `json:"stack"`
			Payload json.RawMessage `json:"payload"`
			Module  string          `json:"module"`
		}{
			ID: DocID{
				Module: log.Module,
				Shard:  log.Shard,
				LogID:  log.ID,
				Stack:  connector.stack,
			}.String(),
			Stack:   connector.stack,
			Payload: log.Payload,
			Module:  log.Module,
		}

		bulk.Add(
			elastic.NewBulkIndexRequest().
				Index(connector.config.Index).
				Id(doc.ID).
				Doc(doc),
		)
	}

	rsp, err := bulk.Do(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query es")
	}

	ret := make([]error, len(logs))
	for index, item := range rsp.Items {
		errorDetails := item["index"].Error
		if errorDetails == nil {
			ret[index] = nil
		} else {
			ret[index] = errors.New(errorDetails.Reason)
		}
	}

	return ret, nil
}

func NewConnector(serviceConfig drivers.ServiceConfig, config Config, logger logging.Logger) (*Connector, error) {
	return &Connector{
		stack:         serviceConfig.Stack,
		serviceConfig: serviceConfig,
		config:        config,
		logger:        logger,
	}, nil
}

var _ drivers.Driver = (*Connector)(nil)

type DocID struct {
	Module string `json:"module"`
	LogID  string `json:"logID"`
	Stack  string `json:"stack"`
	Shard  string `json:"shard,omitempty"`
}

func (docID DocID) String() string {
	rawID, _ := json.Marshal(docID)
	return base64.URLEncoding.EncodeToString(rawID)
}
