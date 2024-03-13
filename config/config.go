package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	cfgFile = flag.String("config", "./config.yaml", "config file path")

	cfg *Config
)

// Config example config
type Config struct {
	Listen string `yaml:"listen"`
	Mode   string `yaml:"mode"`
	Redis  struct {
		Host        string `yaml:"host"`
		Password    string `yaml:"password"`
		Database    int    `yaml:"database"`
		MaxActive   int    `yaml:"maxActive"`
		MaxIdle     int    `yaml:"maxIdle"`
		IdleTimeout int    `yaml:"idleTimeout"`
	} `yaml:"redis"`
	Mysql struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
	JwtSecret        string `yaml:"jwtsecret"`
	*TecentSMSConfig `yaml:"TecentSMSConfig"`
}

type TecentSMSConfig struct {
	SecretId                string `yaml:"SecretId"`
	SecretKey               string `yaml:"SecretKey"`
	SmsSdkAppId             string `yaml:"SmsSdkAppId"`
	LoginTemplateId         string `yaml:"LoginTemplateId"`
	RegisterTemplateId      string `yaml:"RegisterTemplateId"`
	REDIS_KEY_LOGIN_CODE    string `yaml:"REDIS_KEY_LOGIN_CODE"`
	REDIS_KEY_REGISTER_CODE string `yaml:"REDIS_KEY_REGISTER_CODE"`
	ExpirationMinutes       int    `yaml:"ExpirationMinutes"`
}

// GetConfig 获取配置
func GetConfig() *Config {
	if cfg != nil {
		return cfg
	}
	bytes, err := os.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	cfgData := &Config{}
	err = yaml.Unmarshal(bytes, cfgData)
	if err != nil {
		panic(err)
	}
	cfg = cfgData
	return cfg
}
