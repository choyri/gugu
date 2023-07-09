package logx

import (
	"golang.org/x/exp/slog"
)

func NewPkg(pkg string, mod ...string) *Logger {
	attrs := []slog.Attr{
		slog.String(KeyAppPkg, pkg),
	}

	if len(mod) > 0 && mod[0] != "" {
		attrs = append(attrs, slog.String(KeyAppPkgMod, mod[0]))
	}

	return With(attrs...)
}
