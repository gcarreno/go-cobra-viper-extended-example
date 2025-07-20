package config

import (
	"github.com/spf13/viper"
)

// Always define constants early, in order to make changes quick and depend
// on syntax completion to avoid errors
const (
	cLogLevel   = "info"
	cAdminEmail = "webmaster@mysite.com"

	ViperLogLevel   = "log_level"
	ViperAdminEmail = "admin_email"

	ViperWebAddress = "web.address"
	ViperWebPort    = "web.port"

	ViperAPIAddress = "api.address"
	ViperAPIPort    = "api.port"
)

type Config struct {
	BaseConfig `mapstructure:",squash"`
	Web        *WebConfig `mapstructure:"web" json:"web"`
	API        *APIConfig `mapstructure:"api" json:"api"`
}

func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Web:        DefaultWebConfig(),
		API:        DefaultAPIConfig(),
	}
}

type BaseConfig struct {
	LogLevel   string `mapstructure:"log_level" json:"log_level"`
	AdminEmail string `mapstructure:"admin_email" json:"admin_email"`
}

func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		LogLevel:   cLogLevel,
		AdminEmail: cAdminEmail,
	}
}

func SetDefaultsToViper() {
	// Get the defaults
	config := DefaultConfig()

	// BaseConfig
	viper.SetDefault(ViperLogLevel, config.LogLevel)
	viper.SetDefault(ViperAdminEmail, config.AdminEmail)

	// WebConfig
	viper.SetDefault(ViperWebAddress, config.Web.Address)
	viper.SetDefault(ViperWebPort, config.Web.Port)

	// APIConfig
	viper.SetDefault(ViperAPIAddress, config.API.Address)
	viper.SetDefault(ViperAPIPort, config.API.Port)
}
