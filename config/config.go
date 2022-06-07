package conf

import (
	"github.com/spf13/viper"
)

var (
	VP *viper.Viper
	Cf Config
)

type ApiConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Version  string `mapstructure:"version"`
}

type DbConfig struct {
	Dsn string `mapstructure:"dsn"`
}

type JwtConfig struct {
	AccSecKey string `mapstructure:"acc_sec_key"`
	RefSecKey string `mapstructure:"ref_sec_key"`
}

type Config struct {
	Api ApiConfig `mapstructure:"api"`
	Db  DbConfig  `mapstructure:"db"`
	Jwt JwtConfig `mapstructure:"jwt"`
}

func LoadConfig() (c Config, err error) {
	VP = viper.New()
	VP.SetConfigName("config")
	VP.SetConfigType("json")
	VP.AddConfigPath("/config.")
	VP.AddConfigPath(".")
	if err = VP.ReadInConfig(); err != nil {
		return Config{}, err
	}
	if err = VP.Unmarshal(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

func init() {
	var err error
	Cf, err = LoadConfig()
	if err != nil {
		panic(err.Error())
	}
}
