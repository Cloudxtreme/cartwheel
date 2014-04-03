package config

type BackendConfig struct {
	Address string `json:"addr"`
	UseSSL  bool   `json:"use_ssl"`
}
