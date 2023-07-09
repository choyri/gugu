package logx

import (
	"os"

	"golang.org/x/exp/slog"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.Logger.LogAttrs(nil, slog.LevelInfo, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	l.Logger.LogAttrs(nil, slog.LevelError, msg, attrs...)
}

func NewLogger(cfg Config, attrs ...slog.Attr) (*Logger, error) {
	var lvl slog.Level

	err := lvl.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, err
	}

	hdlOpt := slog.HandlerOptions{Level: lvl}

	var hdl slog.Handler

	if cfg.Development {
		hdl = slog.NewTextHandler(os.Stderr, &hdlOpt)
	} else {
		hdl = slog.NewJSONHandler(os.Stderr, &hdlOpt)
	}

	if len(attrs) > 0 {
		hdl = hdl.WithAttrs(attrs)
	}

	logger := slog.New(hdl)
	slog.SetDefault(logger)

	return &Logger{
		Logger: logger,
	}, nil
}

func With(attrs ...slog.Attr) *Logger {
	args := make([]any, 0, len(attrs))

	for _, v := range attrs {
		args = append(args, v)
	}

	return &Logger{
		Logger: slog.Default().With(args...),
	}
}
