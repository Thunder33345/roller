package json

type Option func(*JSON)

func WithReadonly() Option {
	return func(j *JSON) {
		j.readOnly = true
	}
}

func WithAllowUnknown() Option {
	return func(j *JSON) {
		j.allowUnknown = true
	}
}

func WithIndent(indent string) Option {
	return func(j *JSON) {
		j.indent = indent
	}
}
