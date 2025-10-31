package internal

import (
	"go_blog/internal/migration"
	"go_blog/internal/router"

	"go.uber.org/fx"
)

var Module = fx.Options(
	migration.Module,
	router.Module,
)
