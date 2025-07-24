package shortener

import "go.uber.org/fx"

func FxModule() fx.Option {
	return fx.Module(
		"shortener",
		fx.Provide(
			NewServer,
			NewGRPCServer,
			NewHTTPProxy,
		),
	)
}
