package clickhouse

import (
	"context"
	"fmt"

	ingester "github.com/formancehq/stack/ee/ingester/internal"
	"github.com/formancehq/stack/ee/ingester/internal/config"
	"github.com/formancehq/stack/ee/ingester/internal/drivers"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/pkg/errors"
)

type Connector struct {
	db            driver.Conn
	serviceConfig drivers.ServiceConfig
	config        Config
	logger        logging.Logger
}

func (c *Connector) Stop(_ context.Context) error {
	return c.db.Close()
}

func (c *Connector) Start(ctx context.Context) error {

	var err error
	c.db, err = OpenDB(c.logger, c.config.DSN, c.serviceConfig.Debug)
	if err != nil {
		return errors.Wrap(err, "opening database")
	}

	// Create the database
	// One database is used for the entire stack
	err = c.db.Exec(ctx, fmt.Sprintf(`create database if not exists "%s"`, c.serviceConfig.Stack))
	if err != nil {
		return errors.Wrap(err, "failed to create database")
	}

	// Create the logs table
	// One table is used for the entire stack
	err = c.db.Exec(ctx, createLogsTable)
	if err != nil {
		return errors.Wrap(err, "failed to create logs table")
	}

	return nil
}

func (c *Connector) Accept(ctx context.Context, logs ...ingester.LogWithModule) ([]error, error) {

	batch, err := c.db.PrepareBatch(ctx, "insert into logs(module, shard, id, type, date, data)")
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare batch")
	}

	for _, log := range logs {
		if err := batch.Append(
			log.Module,
			log.Shard,
			log.ID,
			log.Type,
			log.Date,
			string(log.Payload),
		); err != nil {
			return nil, errors.Wrap(err, "appending item to the batch")
		}
	}

	return make([]error, len(logs)), errors.Wrap(batch.Send(), "failed to commit transaction")
}

func NewConnector(serviceConfig drivers.ServiceConfig, config Config, logger logging.Logger) (*Connector, error) {
	return &Connector{
		serviceConfig: serviceConfig,
		config:        config,
		logger:        logger,
	}, nil
}

var _ drivers.Driver = (*Connector)(nil)

type Config struct {
	DSN string `json:"dsn"`
}

func (cfg Config) Validate() error {
	if cfg.DSN == "" {
		return errors.New("dsn is required")
	}

	return nil
}

var _ config.Validator = (*Config)(nil)

const createLogsTable = `
	create table if not exists logs (
		module String,
		shard String,
		id              String,
		type            String,
		date            DateTime,
		data            String
	) 
	engine = ReplacingMergeTree
	partition by module
	primary key (module, shard, id)
`

func OpenDB(logger logging.Logger, dsn string, debug bool) (driver.Conn, error) {
	// Open database connection
	options, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "parsing dsn")
	}
	if debug {
		options.Debug = true
		options.Debugf = logger.Debugf
	}

	db, err := clickhouse.Open(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open db")
	}

	return db, nil
}
