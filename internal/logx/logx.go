package logx

import (
	"os"

	"golang.org/x/exp/slog"
)

func Init(conf Config, attrs ...slog.Attr) (*slog.Logger, error) {
	var lvl slog.Level

	err := lvl.UnmarshalText([]byte(conf.Level))
	if err != nil {
		return nil, err
	}

	var (
		hdlOpt = slog.HandlerOptions{Level: lvl}
		hdl    slog.Handler
	)

	if conf.Development {
		hdl = slog.NewTextHandler(os.Stderr, &hdlOpt)
	} else {
		hdl = slog.NewJSONHandler(os.Stderr, &hdlOpt)
	}

	if len(attrs) > 0 {
		hdl = hdl.WithAttrs(attrs)
	}

	logger := slog.New(hdl)
	slog.SetDefault(logger)

	return logger, nil
}

func NewPkg(pkg string, mod ...string) *slog.Logger {
	attrs := []any{
		slog.String(KeyAppPkg, pkg),
	}

	if len(mod) > 0 && mod[0] != "" {
		attrs = append(attrs, slog.String(KeyAppPkgMod, mod[0]))
	}

	return slog.With(attrs...)
}
