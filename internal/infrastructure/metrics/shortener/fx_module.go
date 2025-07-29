package metrics

import "go.uber.org/fx"

func FxModule() fx.Option {
	return fx.Module(
		"metrics",
		fx.Provide(New),
	)
}
