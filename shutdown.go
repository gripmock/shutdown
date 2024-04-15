package shutdown

import (
	"context"
)

type Fn func(context.Context) error

type Logger interface {
	Err(format error)
}

type Shutdown struct {
	fn     []Fn
	logger Logger
}

func New(logger Logger) *Shutdown {
	return &Shutdown{fn: []Fn{}, logger: logger}
}

func (s *Shutdown) Add(fns ...Fn) {
	s.fn = append(s.fn, fns...)
}

func (s *Shutdown) Do(ctx context.Context) {
	for i := len(s.fn) - 1; i >= 0; i-- {
		if err := s.fn[i](ctx); s.logger != nil && err != nil {
			s.logger.Err(err)
		}
	}
}
