package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type config struct {
	Server server `yaml:"server"`
	Db     db     `yaml:"db"`
	Logger logger `yaml:"logger"`
	Redis  redis  `yaml:"redis"`
}

type server struct {
	Port      string `yaml:"port"`
	Name      string `yaml:"name"`
	Mode      string `yaml:"mode"`
	JWTSecret string `yaml:"jwt_secret"`
	JWTExpire int64  `yaml:"jwt_expire"`
}

type db struct {
	Host            string `yaml:"host"`
	Port            int32  `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

type logger struct {
	Level  string `yaml:"level"`
	Path   string `yaml:"path"`
	MaxAge int64  `yaml:"max_age"`
}

type redis struct {
	Host     string `yaml:"host"`
	Port     int32  `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

var Config *config

func init() {
	// 读取配置文件
	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err)
	}

	Config = &config{}
	err = yaml.Unmarshal(file, Config)
	if err != nil {
		panic(err)
	}
}
