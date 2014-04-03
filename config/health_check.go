package config

type HealthCheckConfig struct {
	CheckType      string            `json:"type"`
	CheckArguments map[string]string `json:"args"`
}
