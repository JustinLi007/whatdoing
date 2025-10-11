package configs

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Builder struct {
	mtx             *sync.RWMutex
	expectedCliOpts []string
	expectedEnvOpts []string
}

func NewBuilder() *Builder {
	b := &Builder{
		mtx:             &sync.RWMutex{},
		expectedCliOpts: make([]string, 0),
		expectedEnvOpts: make([]string, 0),
	}
	return b
}

func (b *Builder) Cli(opt string) *Builder {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.expectedCliOpts == nil {
		return b
	}
	b.expectedCliOpts = append(b.expectedCliOpts, opt)
	return b
}

func (b *Builder) Env(opt string) *Builder {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	if b.expectedEnvOpts == nil {
		return b
	}
	b.expectedEnvOpts = append(b.expectedEnvOpts, opt)
	return b
}

func (b *Builder) Build() *Config {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	cliTable := make(map[string]string)
	envTable := make(map[string]string)
	for _, v := range b.expectedCliOpts {
		cliTable[v] = ""
	}
	for _, v := range b.expectedEnvOpts {
		envTable[v] = ""
	}
	return newConfig(cliTable, envTable)
}

func (b *Builder) String() string {
	b.mtx.RLock()
	defer b.mtx.RUnlock()

	var buf bytes.Buffer
	buf.WriteString("ENV Options\n")
	for _, v := range b.expectedEnvOpts {
		buf.WriteString(fmt.Sprintf("%s\n", v))
	}
	buf.WriteString("CLI Options\n")
	for _, v := range b.expectedCliOpts {
		buf.WriteString(fmt.Sprintf("%s\n", v))
	}

	return buf.String()
}

type Config struct {
	mtx     *sync.RWMutex
	cliOpts map[string]string
	envOpts map[string]string
}

func newConfig(cliOpts, envOpts map[string]string) *Config {
	c := &Config{
		mtx:     &sync.RWMutex{},
		cliOpts: make(map[string]string),
		envOpts: make(map[string]string),
	}
	if cliOpts != nil {
		c.cliOpts = cliOpts
	}
	if envOpts != nil {
		c.envOpts = envOpts
	}
	return c
}

func (c *Config) Parse() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.parseEnv()
	c.parseCli()
}

func (c *Config) parseCli() {
	args := os.Args
	i := 0
	j := 1
	for j < len(args) {
		flag := args[i]
		flag = strings.TrimPrefix(flag, "--")
		val := args[j]
		_, ok := c.cliOpts[flag]
		if ok {
			c.cliOpts[flag] = val
		}
		i++
		j++
	}
}

func (c *Config) parseEnv() {
	temp := make([]string, 0)
	for k := range c.envOpts {
		temp = append(temp, k)
	}
	for _, v := range temp {
		c.envOpts[v] = os.Getenv(v)
	}
}

func (c *Config) Get(opt string) string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	var val string
	var ok bool
	if val, ok = c.getCli(opt); ok {
		return val
	} else if val, ok = c.getEnv(opt); ok {
		return val
	}
	return ""
}

func (c *Config) getCli(opt string) (string, bool) {
	if c.cliOpts == nil {
		return "", false
	}
	val, ok := c.cliOpts[opt]
	return val, ok
}

func (c *Config) getEnv(opt string) (string, bool) {
	if c.envOpts == nil {
		return "", false
	}
	val, ok := c.envOpts[opt]
	return val, ok
}
