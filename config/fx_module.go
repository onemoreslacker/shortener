package config

import (
	"flag"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(func() (*Config, error) {
			path := flag.String("path", "", "path to the config file")
			flag.Parse()
			return New(*path)
		}),
	)
}
