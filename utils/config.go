package utils

import (
	"go-web/appError"
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

type PoolConfig struct {
	Size int `yaml:"size"`
}

type JWTConfig struct {
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

type AppConfig struct {
	MySQL MySQLConfig `yaml:"mysql"`
	Pool  PoolConfig  `yaml:"pool"`
	Redis RedisConfig `yaml:"redis"`
	JWT   JWTConfig   `yaml:"jwt"`
}

var Conf AppConfig

func InitConfig() {
	yamlFile, err := os.ReadFile("config/app.yml")
	if err != nil {
		appError.NewAppError(508, "读取配置文件失败", err)
	}
	viper.AutomaticEnv()

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		appError.NewAppError(508, "解析配置文件失败", err)
	}
	viper.SetConfigFile("config/app.yml")
	log.Println("✅ 配置加载成功")
}
