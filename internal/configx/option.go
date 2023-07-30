package configx

import (
	"github.com/choyri/gugu/internal/optionx"
	"github.com/knadh/koanf/v2"
)

type Option struct {
	customParser koanf.Parser
	defaultData  []byte
	envPrefix    string
	filepath     string
}

func WithCustomParser(pa koanf.Parser) optionx.Option[Option] {
	return optionx.OptionFunc[Option](func(opt *Option) {
		opt.customParser = pa
	})
}

func WithDefaultData(data []byte) optionx.Option[Option] {
	return optionx.OptionFunc[Option](func(opt *Option) {
		opt.defaultData = data
	})
}

func WithEnvPrefix(prefix string) optionx.Option[Option] {
	return optionx.OptionFunc[Option](func(opt *Option) {
		opt.envPrefix = prefix
	})
}

func WithFile(filepath string) optionx.Option[Option] {
	return optionx.OptionFunc[Option](func(opt *Option) {
		opt.filepath = filepath
	})
}
