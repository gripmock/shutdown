package shutdown_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gripmock/shutdown"
)

var ErrSame = errors.New("some error")
var ErrSame1 = errors.New("some error 1")
var ErrSame2 = errors.New("some error 2")
var ErrSame3 = errors.New("some error 3")

type logger struct {
	errors []error
}

func (l *logger) Err(err error) {
	l.errors = append(l.errors, err)
}

func TestShutdown_LoggerNil(t *testing.T) {
	s := shutdown.New(nil)
	s.Add(func(ctx context.Context) error {
		return ErrSame
	})

	s.Do(context.Background())

	require.True(t, true)
}

func TestShutdown_Stack_ErrorAll(t *testing.T) {
	l := &logger{}
	s := shutdown.New(l)
	s.Add(
		func(ctx context.Context) error { return ErrSame1 },
		func(ctx context.Context) error { return ErrSame2 },
		func(ctx context.Context) error { return ErrSame3 },
	)

	s.Do(context.Background())

	require.ErrorIs(t, l.errors[0], ErrSame3)
	require.ErrorIs(t, l.errors[1], ErrSame2)
	require.ErrorIs(t, l.errors[2], ErrSame1)
}

func TestShutdown_Stack_Error(t *testing.T) {
	l := &logger{}
	s := shutdown.New(l)
	s.Add(
		func(ctx context.Context) error { return ErrSame1 },
		func(ctx context.Context) error { return nil },
		func(ctx context.Context) error { return ErrSame3 },
	)

	s.Do(context.Background())

	require.ErrorIs(t, l.errors[0], ErrSame3)
	require.ErrorIs(t, l.errors[1], ErrSame1)
}
