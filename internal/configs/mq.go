package configs

type MQ struct {
	Addresses []string `yaml:"addresses"`
	ClientID  string   `yaml:"client_id"`
}
