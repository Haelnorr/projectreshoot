package tests

import (
	"os"
	"projectreshoot/config"

	"github.com/pkg/errors"
)

func TestConfig() (*config.Config, error) {
	os.Setenv("SECRET_KEY", ".")
	cfg, err := config.GetConfig(map[string]string{})
	if err != nil {
		return nil, errors.Wrap(err, "config.GetConfig")
	}
	return cfg, nil
}
