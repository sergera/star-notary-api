package conf

import (
	"log"

	"github.com/gurkankaymak/hocon"
)

var Port string
var LogPath string
var DBHost string
var DBPort string
var DBName string
var DBUser string
var DBPassword string
var CORSAllowedURLs string

var config *hocon.Config

func Setup() {
	parseHOCONConfigFile()
	setPort()
	setLogPath()
	setDBHost()
	setDBPort()
	setDBName()
	setDBUser()
	setDBPassword()
	setCORSAllowedURLs()
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

func setLogPath() {
	LogPath = config.GetString("log-path")
}

func setDBHost() {
	dbHost := config.GetString("db.host")
	if len(dbHost) == 0 {
		log.Panic("database host environment variable not found")
	}

	DBHost = dbHost
}

func setDBPort() {
	dbPort := config.GetString("db.port")
	if len(dbPort) == 0 {
		log.Panic("database port environment variable not found")
	}

	DBPort = dbPort
}

func setDBName() {
	dbName := config.GetString("db.name")

	if len(dbName) == 0 {
		log.Panic("database name environment variable not found")
	}

	DBName = dbName
}

func setDBUser() {
	dbUser := config.GetString("db.user")
	if len(dbUser) == 0 {
		log.Panic("database user environment variable not found")
	}

	DBUser = dbUser
}

func setDBPassword() {
	dbPassword := config.GetString("db.password")
	if len(dbPassword) == 0 {
		log.Panic("database password environment variable not found")
	}

	DBPassword = dbPassword
}

func setCORSAllowedURLs() {
	corsAllowedURLs := config.GetString("cors.urls")
	if len(corsAllowedURLs) == 0 {
		log.Panic("cors allowed urls environment variable not found")
	}

	CORSAllowedURLs = corsAllowedURLs
}
