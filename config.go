package jwttransport

import (
	"strings"
	"io/ioutil"
	"encoding/json"
)

type Config struct {
	Type         string `json:"type"`
	Scope        string `json:"scope"`
	ProjectId    string `json:"project_id"`
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey   string `json:"private_key" binding:"required"`
	ClientEmail  string `json:"client_email" binding:"required"`
	ClientId     string `json:"client_id" binding:"required"`

	// TokenCache allows tokens to be cached for subsequent requests.
	TokenCache 	 Cache
}

type Configurator struct {
	Scope        string
	TokenCache 	 Cache
}

func (c *Configurator) Load(configFile string) (*Config, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	cfg.Scope = c.Scope
	cfg.TokenCache = c.TokenCache

	if cfg.TokenCache == nil {
		cfg.TokenCache = &CacheInMemory{}
	}

	if len(cfg.ProjectId) == 0 && len(cfg.ClientId) > 0 {
		cfg.ProjectId = strings.SplitN(cfg.ClientId, "-", 2)[0]
	}

	return &cfg, nil
}
