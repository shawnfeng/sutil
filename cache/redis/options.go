package redis

// redis Client options
type options struct {
	// fix key  #{namespace.wrapper.key}
	wrapper    string
	// if true no fix key
	noFixKey      bool
	// if true key => #{namespace.wrapper.key} else key => #{namespace.key}
	useWrapper bool
}

type Option interface {
	apply(*options)
}

type wrapperOption string

func (c wrapperOption) apply(opts *options) {
	opts.wrapper = string(c)
}

func WithWrapper(w string) Option {
	return wrapperOption(w)
}

type noFixKeyOption bool

func (c noFixKeyOption) apply(opts *options) {
	opts.noFixKey = bool(c)
}

func WithNoFixKey(n bool) Option {
	return noFixKeyOption(n)
}

type useWrapperOption bool

func (c useWrapperOption) apply(opts *options) {
	opts.useWrapper = bool(c)
}

func WithUseWrapper(n bool) Option {
	return useWrapperOption(n)
}
