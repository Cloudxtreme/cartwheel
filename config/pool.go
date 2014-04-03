package config

type PoolConfig struct {
	Listeners    []ListenerConfig    `json:"listeners"`
	Backends     []BackendConfig     `json:"backends"`
	Balancer     BalancerConfig      `json:"balancer"`
	HealthChecks []HealthCheckConfig `json:"health_checks"`
}
