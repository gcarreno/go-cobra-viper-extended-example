package config

type WebConfig struct {
	Address string `mapstructure:"address" json:"address"`
	Port    int32  `mapstructure:"port" json:"port"`
}

const (
	cWebAddress       = "0.0.0.0"
	cWebPort    int32 = 8080
)

func DefaultWebConfig() *WebConfig {
	return &WebConfig{
		Address: cWebAddress,
		Port:    cWebPort,
	}
}
