package main

import (
	"context"

	"JungleGaming-test/internal/di"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		di.Module(),
		fx.Invoke(func(lc fx.Lifecycle) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error { return nil },
				OnStop: func(context.Context) error { return nil },
			})
		}),
	)
	app.Run()
}
