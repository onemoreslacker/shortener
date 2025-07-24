package mongo

import "go.uber.org/fx"

func FxModule() fx.Option {
	return fx.Module(
		"mongo client",
		fx.Provide(NewClient),
	)
}
