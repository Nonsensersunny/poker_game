package conf

import "fmt"

var (
	DefaultConfig *GlobalConfig
)

type GlobalConfig struct {
	RedisConfig
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%v:%v", r.Host, r.Port)
}

func NewGlobalConfig() (*GlobalConfig, error) {
	return &GlobalConfig{
		RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
	}, nil
}

func init() {
	conf, err := NewGlobalConfig()
	if err != nil {
		panic(err)
	}
	DefaultConfig = conf
}
