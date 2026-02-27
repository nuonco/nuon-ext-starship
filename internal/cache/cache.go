package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Entry struct {
	OrgID       string    `json:"org_id"`
	OrgName     string    `json:"org_name"`
	InstallID   string    `json:"install_id"`
	InstallName string    `json:"install_name"`
	APIURL      string    `json:"api_url"`
	CachedAt    time.Time `json:"cached_at"`
}

func cacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache", "nuon"), nil
}

func cachePath() (string, error) {
	dir, err := cacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "starship.json"), nil
}

func Read() (*Entry, error) {
	p, err := cachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

func Write(e *Entry) error {
	dir, err := cacheDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	e.CachedAt = time.Now()

	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}

	p, err := cachePath()
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, 0o644)
}
