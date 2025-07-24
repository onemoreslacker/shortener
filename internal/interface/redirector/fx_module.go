package redirector

import "go.uber.org/fx"

func FxModule() fx.Option {
	return fx.Module(
		"redirector",
		fx.Provide(
			NewServer,
			NewGRPCServer,
			NewHTTPProxy,
		),
	)
}
