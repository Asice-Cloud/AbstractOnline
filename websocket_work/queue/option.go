package queue

type handlerfunc func(msg Message)

type Options struct {
	topic   string
	handler handlerfunc
}

type Option func(*Options)

func WithTopic(topic string) Option {
	return func(o *Options) {
		o.topic = topic
	}
}

func WithHandler(handler handlerfunc) Option {
	return func(o *Options) {
		o.handler = handler
	}
}
