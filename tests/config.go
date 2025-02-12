package tests

import (
	"os"
	"projectreshoot/config"

	"github.com/pkg/errors"
)

func TestConfig() (*config.Config, error) {
	os.Setenv("TRUSTED_HOST", "127.0.0.1")
	os.Setenv("TURSO_DB_NAME", ".")
	os.Setenv("TURSO_AUTH_TOKEN", ".")
	os.Setenv("SECRET_KEY", ".")
	cfg, err := config.GetConfig(map[string]string{})
	if err != nil {
		return nil, errors.Wrap(err, "config.GetConfig")
	}
	return cfg, nil
}
