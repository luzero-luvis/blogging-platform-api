package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBURL string
	PORT  string
}

func Load() (*Config, error) {
	conf := &Config{
		DBURL: os.Getenv("DB_URL"),
		PORT:  os.Getenv("PORT"),
	}

	required := map[string]string{
		"DBURL": conf.DBURL,
	}

	var missing []string
	for key, val := range required {
		if val == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		for _, key := range missing {
			fmt.Printf("env required %s\n", key)
		}
		return nil, fmt.Errorf("env required %s", strings.Join(missing, " "))
	}
	return conf, nil
}
