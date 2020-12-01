package inttree

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type loggingIntTree struct {
	logger log.Logger
	tree   Tree
}

func (l *loggingIntTree) Find(val int) (foundData int, err error) {
	defer func(from time.Time) {
		_ = l.wrappedLogger(err).Log(
			"method", "Find",
			"value", val,
			"foundData", foundData,
			"err", err,
			"executionTime", time.Since(from),
		)
	}(time.Now())
	return l.tree.Find(val)
}

func (l *loggingIntTree) Append(val int) (err error) {
	defer func(from time.Time) {
		_ = l.wrappedLogger(err).Log(
			"method", "Append",
			"value", val,
			"err", err,
			"executionTime", time.Since(from),
		)
	}(time.Now())
	return l.tree.Append(val)
}

func (l *loggingIntTree) Delete(val int) (err error) {
	defer func(from time.Time) {
		_ = l.wrappedLogger(err).Log(
			"method", "Delete",
			"value", val,
			"err", err,
			"executionTime", time.Since(from),
		)
	}(time.Now())
	return l.tree.Delete(val)
}

func (l *loggingIntTree) wrappedLogger(err error) log.Logger {
	if err != nil {
		return level.Error(l.logger)
	}
	return l.logger
}

func NewLoggingIntTree(logger log.Logger, tree Tree) Tree {
	return &loggingIntTree{
		logger: logger,
		tree:   tree,
	}
}
