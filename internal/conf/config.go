package conf

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	once sync.Once
	cfg  Config
)

// General config of application
type Config struct {
	// Config application generate and send information
	App struct {
		// Count of facts, which are sended to server one-by-one
		FactsGenCount int `env:"FACTS_GEN_COUNT" env-default:"10"`

		// Count of goroutine, which sends to server in parallel
		FactsGorCount int `env:"FACTS_GOR_COUNT" env-default:"1"`
	}

	// Config for server
	// Set up separately6 because Host and Port may be different
	// but action is same
	Server struct {
		Host   string `env:"SERVER_HOST" env-default:"https://development.kpi-drive.ru/"`
		Port   int    `env:"SERVER_PORT" env-default:"80"`
		Action string `env:"SERVER_ACTION" env-default:"_api/facts/save_fact"`

		Token string `env:"SERVER_TOKEN" env-default:"BEARER_TOKEN"`

		Url string
	}

	// Config for work mode of application
	Launch struct {
		// Mode of work.
		// Only logger depends on it.
		// - "debug" - write to console logs
		// - "release" - write to file logs
		AppMode string `env:"APP_MODE" env-default:"debug"`

		// Log level of application
		LogLevel string `env:"LOG_LEVEL" env-default:"debug"`
	}
}

// Get main config
// Read config execute one time
func Get() Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			panic(err)
		}

		parts := make([]string, 3)
		parts[0] = cfg.Server.Host
		parts[2] = cfg.Server.Action

		if cfg.Server.Port != 80 {
			parts[1] = fmt.Sprintf("%d", cfg.Server.Port)
		}

		cfg.Server.Url = strings.Join(parts, "")
	})
	return cfg
}
