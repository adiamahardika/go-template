package config

import (
	"monitoring-service/pkg/database"
	"net/url"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceHost        string           `mapstructure:"service_host" json:"service_host"`
	ServiceEndpointV   string           `mapstructure:"service_endpoint_v" json:"service_endpoint_v"`
	ServiceEnvironment string           `mapstructure:"service_environment" json:"service_environment"`
	ServicePort        string           `mapstructure:"service_port" json:"service_port"`
	Database           DatabasePlatform `mapstructure:"database" json:"database"`
}

func NewConfig() *Config {
	return &Config{
		ServiceHost:        viper.GetString("APP_HOST"),
		ServiceEndpointV:   viper.GetString("APP_ENDPOINT_V"),
		ServiceEnvironment: viper.GetString("APP_ENVIRONMENT"),
		ServicePort:        viper.GetString("APP_PORT"),
		Database:           LoadDatabaseConfig(),
	}
}

func (d *Database) ToArgs(dbType database.DBType, connType database.ConnType, val url.Values) (res *database.Args) {
	res = &database.Args{
		Username:        d.Username,
		Password:        d.Password,
		Host:            d.URL,
		Port:            d.Port,
		Database:        d.Name,
		Schema:          d.Schema,
		MaxIdleConns:    d.MaxIdleConns,
		MaxOpenConns:    d.MaxOpenConns,
		ConnMaxLifetime: d.MaxLifetime,
		Flavor:          d.Flavor,
		Location:        d.Location,
		Timeout:         d.Timeout,

		DBType:   dbType,
		ConnType: connType,
		Values:   val,
	}
	return
}
