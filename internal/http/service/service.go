package service

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
)

// tree implements tree for integer values
type tree interface {
	// Find finds data by value
	Find(val int) (foundData int, err error)
	// Append appends value to tree
	Append(val int) (err error)
	// Delete deletes value from tree
	Delete(val int) (err error)
}

// Service manages operations on tree
type Service interface {
	// Find finds data by value
	Find(ctx context.Context, value int) (val int, err error)
	// Append appends value to tree
	Append(ctx context.Context, value int) (err error)
	// Delete deletes value from tree
	Delete(ctx context.Context, value int) (err error)
}

var _ Service = (*service)(nil)

type service struct {
	t      tree
	logger log.Logger
}

func (s *service) Find(ctx context.Context, value int) (val int, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
		return s.t.Find(value)
	}
}

func (s *service) Append(ctx context.Context, value int) (err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
		return s.t.Append(value)
	}
}

func (s *service) Delete(ctx context.Context, value int) (err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
		return s.t.Delete(value)
	}
}

func NewService(logger log.Logger, initialTree tree) (svc Service, err error) {
	if initialTree == nil {
		err = errors.New("empty initial tree")
		return
	}

	svc = &service{
		t:      initialTree,
		logger: logger,
	}
	return
}
