package interleaving

type IsInclude = func(interface{}, interface{}) bool

type Option = *option
type option struct {
	IsInclude *IsInclude
}

func WithDuplicate(f IsInclude) Option {
	return &option{
		IsInclude: &f,
	}
}

func mergeOptions(opts []Option) Option {
	var opt option
	for _, o := range opts {
		if o == nil {
			continue
		}
		if o.IsInclude != nil {
			opt.IsInclude = o.IsInclude
		}
	}
	return &opt
}
