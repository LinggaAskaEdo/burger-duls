package lib

import "go.uber.org/fx"

// Module exports dependency
var Module = fx.Options(
	fx.Provide(NewRequestHandler),
	fx.Provide(NewEnv),
	fx.Provide(NewLogger),
	fx.Provide(NewDatabase),
	fx.Provide(NewValidator),
)
