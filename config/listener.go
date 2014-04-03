package config

type ListenerConfig struct {
	ListenAddress           string `json:"listen_addr"`
	SSLCertificatePath      string `json:"ssl_cert_path"`
	SSLCertificateChainPath string `json:"ssl_cert_chain_path"`
	SSLKeyPath              string `json:"ssl_key_path"`
}
