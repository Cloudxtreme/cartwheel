package config

type BalancerConfig struct {
	BalancerType      string            `json:"type"`
	BalancerArguments map[string]string `json:"args"`
}
