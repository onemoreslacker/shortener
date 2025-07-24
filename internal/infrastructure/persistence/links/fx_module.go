package links

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module(
		"links repository",
		fx.Provide(
			NewCollection,
			NewRepository,
		),
	)
}
