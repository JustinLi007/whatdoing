package configs

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/JustinLi007/whatdoing/libs/go/utils"
)

const (
	ENV_DEBUG   = "debug"
	ENV_PROD    = "prod"
	APP_SERVICE = "service"
	APP_PUB     = "pub"
)

type Configs struct {
	ModeEnv      string
	ModeApp      string
	ConfigServer *configServer
	ConfigDb     *configDb
}

type configServer struct {
	Port int
}

type configDb struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

var configsInstance *Configs

func NewConfigs(opts ...ConfigOptionFn) *Configs {
	if configsInstance != nil {
		return configsInstance
	}

	appMode, envMode, err := RequireModes()
	utils.RequireNoError(err, "must provide app and env mode via cli or env vars")

	configOptions := &ConfigOption{
		withServer: false,
		withDb:     false,
	}

	if appMode == APP_SERVICE {
		configOptions.withServer = true
	}

	for _, v := range opts {
		v(configOptions)
	}

	c := &Configs{
		ModeEnv:      envMode,
		ModeApp:      appMode,
		ConfigServer: nil,
		ConfigDb:     nil,
	}

	if configOptions.withServer {
		c.ConfigServer = newConfigServer()
	}
	if configOptions.withDb {
		c.ConfigDb = newConfigDb()
	}

	mode := os.Getenv("MODE")
	if mode == ENV_PROD {
		c.ModeEnv = ENV_PROD
	}

	configsInstance = c
	return configsInstance
}

func newConfigServer() *configServer {
	return &configServer{
		Port: 0,
	}
}

func newConfigDb() *configDb {
	return &configDb{
		User:     "",
		Password: "",
		Host:     "",
		Port:     0,
		Database: "",
	}
}

func (c *Configs) LoadEnv() error {
	if c.ModeEnv == ENV_DEBUG {
		return c.loadDebugEnv()
	}

	if c.ConfigServer != nil {
		serverPortStr := os.Getenv("SERVER_PORT")
		if port, err := strconv.Atoi(serverPortStr); err != nil {
			return err
		} else {
			c.ConfigServer.Port = port
		}
	}

	if c.ConfigDb != nil {
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
	}

	return nil
}

func (c *Configs) loadDebugEnv() error {
	if c.ConfigServer != nil {
		serverPortStr := os.Getenv("SERVER_PORT")
		if port, err := strconv.Atoi(serverPortStr); err != nil {
			return err
		} else {
			c.ConfigServer.Port = port
		}
	}

	if c.ConfigDb != nil {
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
	}

	return nil
}

func (c *configDb) PostgresConnStr() string {
	if c == nil {
		return ""
	}

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

func RequireModes() (appMode, envMode string, err error) {
	appCli, envCli, okCli := fromCLI()
	appEnv, envEnv, okEnv := fromEnv()

	if !okCli && !okEnv {
		return "", "", errors.New("missing modes")
	}

	if okCli {
		appMode = appCli
		envMode = envCli
	} else if okEnv {
		appMode = appEnv
		envMode = envEnv
	}

	if appMode != APP_SERVICE && appMode != APP_PUB {
		return "", "", errors.New("missing modes")
	}
	if envMode != ENV_DEBUG && envMode != ENV_PROD {
		return "", "", errors.New("missing modes")
	}

	return appMode, envMode, nil
}

func fromEnv() (appMode, envMode string, ok bool) {
	appMode = os.Getenv("MODE_APP")
	envMode = os.Getenv("MODE_ENV")

	appMode = strings.ToLower(strings.TrimSpace(appMode))
	envMode = strings.ToLower(strings.TrimSpace(envMode))

	if appMode == "" || envMode == "" {
		return "", "", false
	}

	return appMode, envMode, true
}

func fromCLI() (appMode, envMode string, ok bool) {
	args := os.Args
	if len(args) < 3 {
		return "", "", false
	}

	appMode = strings.ToLower(strings.TrimSpace(args[1]))
	envMode = strings.ToLower(strings.TrimSpace(args[2]))

	if appMode == "" || envMode == "" {
		return "", "", false
	}

	return appMode, envMode, true
}

type ConfigOption struct {
	withServer bool
	withDb     bool
}

type ConfigOptionFn func(co *ConfigOption)

func WithServerConfig() ConfigOptionFn {
	return func(co *ConfigOption) {
		co.withServer = true
	}
}

func WithDbConfig() ConfigOptionFn {
	return func(co *ConfigOption) {
		co.withDb = true
	}
}
