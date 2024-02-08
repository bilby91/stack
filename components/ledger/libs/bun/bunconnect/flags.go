package bunconnect

import (
	"context"
	"database/sql/driver"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/lib/pq"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"time"
)

const (
	PostgresURIFlag             = "postgres-uri"
	PostgresAWSEnableIAMFlag    = "postgres-aws-enable-iam"
	PostgresMaxIdleConnsFlag    = "postgres-max-idle-conns"
	PostgresMaxOpenConnsFlag    = "postgres-max-open-conns"
	PostgresConnMaxIdleTimeFlag = "postgres-conn-max-idle-time"
)

func InitFlags(flags *pflag.FlagSet) {
	flags.String(PostgresURIFlag, "", "Postgres URI")
	flags.Bool(PostgresAWSEnableIAMFlag, false, "Enable AWS IAM authentication")
	flags.Int(PostgresMaxIdleConnsFlag, 0, "Max Idle connections")
	flags.Duration(PostgresConnMaxIdleTimeFlag, time.Minute, "Max Idle time for connections")
	flags.Int(PostgresMaxOpenConnsFlag, 20, "Max opened connections")
}

func ConnectionOptionsFromFlags(v *viper.Viper, output io.Writer, debug bool) (*ConnectionOptions, error) {
	var connector func(string) (driver.Connector, error)

	if v.GetBool(PostgresAWSEnableIAMFlag) {
		cfg, err := config.LoadDefaultConfig(context.Background(), iam.LoadOptionFromViper(v))
		if err != nil {
			return nil, err
		}

		connector = func(s string) (driver.Connector, error) {
			return &iamConnector{
				dsn: s,
				driver: &iamDriver{
					awsConfig: cfg,
				},
			}, nil
		}
	} else {
		connector = func(dsn string) (driver.Connector, error) {
			return pq.NewConnector(dsn)
		}
	}
	return &ConnectionOptions{
		DatabaseSourceName: v.GetString(PostgresURIFlag),
		Debug:              debug,
		Writer:             output,
		MaxIdleConns:       v.GetInt(PostgresMaxIdleConnsFlag),
		ConnMaxIdleTime:    v.GetDuration(PostgresConnMaxIdleTimeFlag),
		MaxOpenConns:       v.GetInt(PostgresMaxOpenConnsFlag),
		Connector:          connector,
	}, nil
}
