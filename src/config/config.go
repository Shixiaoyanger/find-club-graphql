package config

import (
	"io/ioutil"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
)

const ConfigFilePrefix = "D:/GoProject/find-club-graphql/src/config/"

type Config struct {
	AppInfo appInfo `json:"Appinfo"`
	Log     logconf `json:"Log"`
	Redis   redis   `json:"Redis"`
	Mysql   mysql   `json:"Mysql"`
	MongoDB mongodb `json:"MongoDB"`
}

type appInfo struct {
	Env  string `json:"Env"`
	Addr string `json:"Addr"`
}

type logconf struct {
	LogBasePath string `json:"LogBasePath"`
	LogFileName string `josn:"LogFileName"`
}
type redis struct {
	Host string
	Port string
	PW   string
}
type mysql struct {
	Host   string
	Port   string
	User   string
	DBname string
	PW     string
}
type mongodb struct {
	DriverName  string `json:"DriverName"`
	Host        string `json:"Host"`
	Port        string `json:"Port"`
	DBName      string `json:"DBName"`
	User        string `json:"User"`
	PW          string `json:"PW"`
	AdminDBName string `json:"AdminDBName"`
}

var Conf *Config

func init() {
	log.Println("Begin init all configs")
	initConf()
	log.Println("over init all configs ")
}

func initConf() {
	log.Println("Begin init config")
	Conf = &Config{}
	fileName := "default.json"
	if v, ok := os.LookupEnv("ENV"); ok {
		fileName = v + ".json"
	}

	filePrefix := ConfigFilePrefix

	if v, ok := os.LookupEnv("CONFIG_PATH_PREFIX"); ok {
		filePrefix = v
	}
	log.Println("ConfigPrefix is ", filePrefix)

	data, err := ioutil.ReadFile(filePrefix + fileName)
	if err != nil {
		log.Println("config-initConf: read default.json error", filePrefix+fileName)
		return
	}
	err = jsoniter.Unmarshal(data, Conf)
	if err != nil {
		log.Println("config-initConf: unmarshal default.json error")
		log.Panic(err)
		return
	}
	if v, ok := os.LookupEnv("ENV"); ok {
		fileName = v + ".json"
		log.Println("Init " + v + " configs successfully!")
	} else {
		log.Println("Init default configs successfully!")
	}

}
