package conf

import (
	"log"
	"sync"

	"github.com/gurkankaymak/hocon"
)

var once sync.Once
var instance *conf
var err error

type conf struct {
	hocon           *hocon.Config
	Port            string
	LogPath         string
	DBHost          string
	DBPort          string
	DBName          string
	DBUser          string
	DBPassword      string
	CORSAllowedURLs string
}

func ConfSingleton() (*conf, error) {
	once.Do(func() {
		var c *conf = &conf{}
		var e = c.setup()
		instance = c
		err = e
	})

	if instance != nil {
		return instance, nil
	}

	return nil, err
}

func (c *conf) setup() error {
	err := c.parseHOCONConfigFile()
	if err != nil {
		log.Printf("error while parsing configuration file: %s", err)
		return err
	}
	c.readPort()
	c.readLogPath()
	c.readDBHost()
	c.readDBPort()
	c.readDBName()
	c.readDBUser()
	c.readDBPassword()
	c.readCORSAllowedURLs()
	return nil
}

func (c *conf) parseHOCONConfigFile() error {
	hocon, err := hocon.ParseResource("application.conf")
	if err != nil {
		return err
	}

	log.Printf("configurations: %+v", *hocon)
	c.hocon = hocon
	return nil
}

func (c *conf) readPort() {
	port := c.hocon.GetString("host.port")
	if len(port) == 0 {
		log.Panic("port environment variable not found")
	}

	c.Port = port
}

func (c *conf) readLogPath() {
	c.LogPath = c.hocon.GetString("log.path")
}

func (c *conf) readDBHost() {
	dbHost := c.hocon.GetString("db.host")
	if len(dbHost) == 0 {
		log.Panic("database host environment variable not found")
	}

	c.DBHost = dbHost
}

func (c *conf) readDBPort() {
	dbPort := c.hocon.GetString("db.port")
	if len(dbPort) == 0 {
		log.Panic("database port environment variable not found")
	}

	c.DBPort = dbPort
}

func (c *conf) readDBName() {
	dbName := c.hocon.GetString("db.name")

	if len(dbName) == 0 {
		log.Panic("database name environment variable not found")
	}

	c.DBName = dbName
}

func (c *conf) readDBUser() {
	dbUser := c.hocon.GetString("db.user")
	if len(dbUser) == 0 {
		log.Panic("database user environment variable not found")
	}

	c.DBUser = dbUser
}

func (c *conf) readDBPassword() {
	dbPassword := c.hocon.GetString("db.password")
	if len(dbPassword) == 0 {
		log.Panic("database password environment variable not found")
	}

	c.DBPassword = dbPassword
}

func (c *conf) readCORSAllowedURLs() {
	corsAllowedURLs := c.hocon.GetString("cors.urls")
	if len(corsAllowedURLs) == 0 {
		log.Panic("cors allowed urls environment variable not found")
	}

	c.CORSAllowedURLs = corsAllowedURLs
}
