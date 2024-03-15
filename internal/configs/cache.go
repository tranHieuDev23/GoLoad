package configs

type CacheType string

const (
	CacheTypeInMemory CacheType = "in_memory"
	CacheTypeRedis    CacheType = "redis"
)

type Cache struct {
	Type     CacheType `yaml:"type"`
	Address  string    `yaml:"address"`
	Username string    `yaml:"username"`
	Password string    `yaml:"password"`
}
