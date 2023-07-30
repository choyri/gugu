package configx

import (
	"fmt"
	"os"
	"strings"

	"github.com/choyri/gugu/internal/optionx"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"golang.org/x/exp/slog"
)

const defaultFileName = "config.toml"

func Init(receiver any, opts ...optionx.Option[Option]) error {
	k := koanf.New(".")

	var (
		opt              = optionx.New(opts)
		pa  koanf.Parser = toml.Parser()
	)

	if opt.customParser != nil {
		pa = opt.customParser
	}

	if len(opt.defaultData) > 0 {
		err := k.Load(rawbytes.Provider(opt.defaultData), pa)
		if err != nil {
			return fmt.Errorf("configx: load default data: %w", err)
		}
	}

	if fi, err := os.Stat(defaultFileName); err == nil {
		err = k.Load(file.Provider(defaultFileName), pa)
		if err != nil {
			slog.Warn("failed to load default file", slog.String("filename", fi.Name()))
		}
	}

	if opt.filepath != "" {
		err := k.Load(file.Provider(opt.filepath), pa)
		if err != nil {
			return fmt.Errorf("configx: load file data: %w", err)
		}
	}

	if opt.envPrefix != "" {
		prefix := opt.envPrefix + "_"
		err := k.Load(env.Provider(prefix, ".", func(s string) string {
			return strings.Replace(strings.ToLower(
				strings.TrimPrefix(s, prefix)), "_", ".", -1)
		}), nil)
		if err != nil {
			return fmt.Errorf("configx: load env data: %w", err)
		}
	}

	err := k.Unmarshal("", receiver)
	if err != nil {
		return fmt.Errorf("configx: unmarshal: %w", err)
	}

	return nil
}
