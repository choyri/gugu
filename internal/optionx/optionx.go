package optionx

type Kind interface {
	any
}

type Option[T Kind] interface {
	Apply(*T)
}

type OptionFunc[T Kind] func(*T)

func (f OptionFunc[T]) Apply(t *T) { f(t) }

func New[T Kind](opts []Option[T], defaults ...Option[T]) *T {
	return Apply(nil, opts, defaults...)
}

func Apply[T Kind](opt *T, opts []Option[T], defaults ...Option[T]) *T {
	if opt == nil {
		opt = new(T)
	}

	for _, o := range defaults {
		o.Apply(opt)
	}

	for _, o := range opts {
		o.Apply(opt)
	}

	return opt
}
