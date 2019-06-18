package util

import (
	"github.com/labstack/gommon/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	// 编译配置
	Build struct {
		Path   string   `yaml:"path"`
		GitUrl string   `yaml:"gitUrl"`
		Lib    []string `yaml:"lib"`
		SSH    SSH      `yaml:"ssh"`
	}

	// 部署配置
	Deploy struct {
		SSH      SSH    `yaml:"ssh"`
		MainFunc string `yaml:"mainFunc"`
	}
}

type SSH struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	IdRsaPath string `yaml:"idRsaPath"`
	UserName  string `yaml:"userName"`
}

var config *Config

func GetConfig() *Config {
	return config
}

/**
 * 加载指定配置文件，读取配置
 */
func LoadConfig(filePath string) {
	// 1. 加载配置文件
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Errorf("cannot load config %s, err=%s", filePath, err.Error())
		panic(err.Error())
	}

	// 2. 读取配置文件
	var c Config
	err = yaml.Unmarshal(content, &c)
	if err != nil {
		log.Errorf("cannot parse config %s, err=%s", filePath, err.Error())
		panic(err.Error())
	}

	config = &c
}

func init() {
	//LoadConfig("conf.yaml")
}
