package configs

import (
	"fmt"
	"os"
	"strconv"
)

const (
	CONFIG_MODE_DEBUG = "DEBUG"
	CONFIG_MODE_PROD  = "PROD"
)

type ConfigMode string

type Configs struct {
	Mode         ConfigMode
	ConfigServer *ConfigServer
	ConfigDb     *ConfigDb
}

type ConfigServer struct {
	Port int
	Iss  string
	Aud  string
}

type ConfigDb struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

func NewConfigs() *Configs {
	c := &Configs{
		Mode:         CONFIG_MODE_DEBUG,
		ConfigServer: newConfigServer(),
		ConfigDb:     newConfigDb(),
	}

	mode := os.Getenv("MODE")
	if mode == CONFIG_MODE_PROD {
		c.Mode = CONFIG_MODE_PROD
	}

	return c
}

func newConfigServer() *ConfigServer {
	return &ConfigServer{
		Port: 0,
		Iss:  "",
	}
}

func newConfigDb() *ConfigDb {
	return &ConfigDb{
		User:     "",
		Password: "",
		Host:     "",
		Port:     0,
		Database: "",
	}
}

func (c *Configs) LoadEnv() error {
	if c.Mode == CONFIG_MODE_DEBUG {
		return c.loadDebugEnv()
	}

	serverPortStr := os.Getenv("SERVER_PORT")
	if port, err := strconv.Atoi(serverPortStr); err != nil {
		return err
	} else {
		c.ConfigServer.Port = port
	}

	c.ConfigServer.Iss = os.Getenv("JWT_ISSUER")
	c.ConfigServer.Aud = os.Getenv("JWT_AUDIENCE")

	c.ConfigDb.User = os.Getenv("DB_USER")
	c.ConfigDb.Password = os.Getenv("DB_PASSWORD")
	c.ConfigDb.Host = os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	if port, err := strconv.Atoi(dbPortStr); err != nil {
		return err
	} else {
		c.ConfigDb.Port = port
	}
	c.ConfigDb.Database = os.Getenv("DB_NAME")
	return nil
}

func (c *Configs) loadDebugEnv() error {
	serverPortStr := os.Getenv("SERVER_PORT")
	if port, err := strconv.Atoi(serverPortStr); err != nil {
		return err
	} else {
		c.ConfigServer.Port = port
	}

	c.ConfigServer.Iss = os.Getenv("JWT_ISSUER")
	c.ConfigServer.Aud = os.Getenv("JWT_AUDIENCE")

	c.ConfigDb.User = os.Getenv("DEBUG_DB_USER")
	c.ConfigDb.Password = os.Getenv("DEBUG_DB_PASSWORD")
	c.ConfigDb.Host = os.Getenv("DEBUG_DB_HOST")
	dbPortStr := os.Getenv("DEBUG_DB_PORT")
	if port, err := strconv.Atoi(dbPortStr); err != nil {
		return err
	} else {
		c.ConfigDb.Port = port
	}
	c.ConfigDb.Database = os.Getenv("DEBUG_DB_NAME")
	return nil
}

func (c *ConfigDb) PostgresConnStr() string {
	if c.User == "" {
		return ""
	}
	if c.Password == "" {
		return ""
	}
	if c.Host == "" {
		return ""
	}
	if c.Port == 0 {
		return ""
	}
	if c.Database == "" {
		return ""
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}
