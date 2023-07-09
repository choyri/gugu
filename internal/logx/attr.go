package logx

import (
	"golang.org/x/exp/slog"
)

const (
	KeyAppPkg    = "app.pkg"
	KeyAppPkgMod = "app.pkg_mod"

	KeyError = "error"
)

func Error(err error) slog.Attr {
	return slog.String(KeyError, err.Error())
}
