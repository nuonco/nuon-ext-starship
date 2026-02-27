package config

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type NuonConfig struct {
	APIToken  string `yaml:"api_token"`
	APIURL    string `yaml:"api_url"`
	OrgID     string `yaml:"org_id"`
	AppID     string `yaml:"app_id"`
	InstallID string `yaml:"install_id"`
}

func Load() (*NuonConfig, error) {
	cfg := &NuonConfig{}

	cfgPath := os.Getenv("NUON_CONFIG_FILE")
	if cfgPath == "" {
		cfgPath = "~/.nuon"
	}

	// Expand ~ to home directory.
	if strings.HasPrefix(cfgPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cfgPath = filepath.Join(home, cfgPath[2:])
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Environment variables override file values.
	if v := os.Getenv("NUON_API_TOKEN"); v != "" {
		cfg.APIToken = v
	}
	if v := os.Getenv("NUON_API_URL"); v != "" {
		cfg.APIURL = v
	}
	if v := os.Getenv("NUON_ORG_ID"); v != "" {
		cfg.OrgID = v
	}
	if v := os.Getenv("NUON_APP_ID"); v != "" {
		cfg.AppID = v
	}
	if v := os.Getenv("NUON_INSTALL_ID"); v != "" {
		cfg.InstallID = v
	}

	if cfg.APIURL == "" {
		cfg.APIURL = "https://api.nuon.co"
	}

	return cfg, nil
}

// EnvName extracts the environment name from the API URL.
//
//	https://api.nuon.co       → prod
//	https://api.stage.nuon.co → stage
//	http://localhost:8081      → local
func (c *NuonConfig) EnvName() string {
	u, err := url.Parse(c.APIURL)
	if err != nil {
		return "unknown"
	}

	host := u.Hostname()
	if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1") {
		return "local"
	}

	// api.nuon.co → prod, api.stage.nuon.co → stage
	parts := strings.Split(host, ".")
	if len(parts) == 3 {
		// api.nuon.co
		return "prod"
	}
	if len(parts) >= 4 {
		// api.stage.nuon.co → parts[1] = "stage"
		return parts[1]
	}

	return "unknown"
}
