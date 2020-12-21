package config

type KubernetesConfig struct {
	Type string  `mapstructure:"type"`
	Kubeconfig string  `mapstructure:"config"`
}
