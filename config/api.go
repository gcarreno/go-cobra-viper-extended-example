package config

type APIConfig struct {
	Address string `mapstructure:"address" json:"address"`
	Port    int32  `mapstructure:"port" json:"port"`
}

const (
	cAPIAddress       = "127.0.0.1"
	cAPIPort    int32 = 80
)

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		Address: cAPIAddress,
		Port:    cAPIPort,
	}
}
