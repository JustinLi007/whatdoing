package configs

import (
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
}

type ConfigServer struct {
	Port     int
	Issuer   string
	Audience string
	JwkUrl   string
}

func NewConfigs() *Configs {
	c := &Configs{
		Mode:         CONFIG_MODE_DEBUG,
		ConfigServer: newConfigServer(),
	}

	mode := os.Getenv("MODE")
	if mode == CONFIG_MODE_PROD {
		c.Mode = CONFIG_MODE_PROD
	}

	return c
}

func newConfigServer() *ConfigServer {
	return &ConfigServer{
		Port:     0,
		Issuer:   "",
		Audience: "",
		JwkUrl:   "",
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

	c.ConfigServer.JwkUrl = os.Getenv("JWK_URL")
	c.ConfigServer.Issuer = os.Getenv("JWT_ISSUER")
	c.ConfigServer.Audience = os.Getenv("JWT_AUDIENCE")

	return nil
}

func (c *Configs) loadDebugEnv() error {
	serverPortStr := os.Getenv("SERVER_PORT")
	if port, err := strconv.Atoi(serverPortStr); err != nil {
		return err
	} else {
		c.ConfigServer.Port = port
	}

	c.ConfigServer.JwkUrl = os.Getenv("JWK_URL")
	c.ConfigServer.Issuer = os.Getenv("JWT_ISSUER")
	c.ConfigServer.Audience = os.Getenv("JWT_AUDIENCE")

	return nil
}
