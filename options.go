package maeve

type Option func(*Options)

type Options struct {
	StorePath  string
	ConfigPath string
}

func WithStorePath(path string) Option {
	return func(options *Options) {
		options.StorePath = path
	}
}

func WithConfigPath(path string) Option {
	return func(options *Options) {
		options.ConfigPath = path
	}
}
