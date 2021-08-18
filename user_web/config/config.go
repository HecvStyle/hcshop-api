package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type JWTConfig struct {
	SignKey string `mapstructure:"key"`
}

type AliSmsConfig struct {
	ApiKey    string `mapstructure:"key"`
	ApiSecret string `mapstructure:"secret"`
	Expire    int    `mapstructure:"expire"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JWTInfo     JWTConfig     `mapstructure:"jwt"`
	AliSmsInfo  AliSmsConfig  `mapstructure:"sm"`
	RedisInfo   RedisConfig   `mapstructure:"redis"`
}
