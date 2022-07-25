package conf

import (
	"log"

	"github.com/gurkankaymak/hocon"
)

var Port string

var config *hocon.Config

func Setup() {
	parseHOCONConfigFile()
	setPort()
}

func parseHOCONConfigFile() {
	conf, err := hocon.ParseResource("application.conf")
	if err != nil {
		log.Panic("error while parsing configuration file: ", err)
	}

	log.Printf("all configuration: %+v", *conf)

	config = conf
}

func setPort() {
	port := config.GetString("host.port")
	if len(port) == 0 {
		log.Panic("port environment variable not found")
	}

	Port = port
}
