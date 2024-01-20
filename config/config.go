package config

import (
	"os"
	"wechat_robot/logrus"

	"gopkg.in/yaml.v2"
)

var serverConf ServerConfig

func GetMySQL() MySQLConf {
	return serverConf.MySQL
}

func GetRedis() RedisConf {
	return serverConf.Redis
}

func GetBaidu() BaiduAppConf {
	return serverConf.Baidu
}

func GetChatGPT() ChatGPTApiConf {
	return serverConf.ChatGPT
}

func Init() {
	content, err := os.ReadFile("./conf/deploy.yml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(content, &serverConf); err != nil {
		panic(err)
	}

	logrus.GetLogger().Debugf("conf: %+v", serverConf)
}

type ServerConfig struct {
	MySQL   MySQLConf      `yaml:"mysql"`
	Redis   RedisConf      `yaml:"redis"`
	Baidu   BaiduAppConf   `yaml:"baidu"`
	ChatGPT ChatGPTApiConf `yaml:"chatgpt"`
}

type MySQLConf struct {
	DBName   string `yaml:"db_name"`
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type RedisConf struct {
	IP       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type BaiduAppConf struct {
	AppKey    string `yaml:"app_key"`
	AppSecret string `yaml:"app_secret"`
}

type ChatGPTApiConf struct {
	ApiKey string `yaml:"api_key"`
}
