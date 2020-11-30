package inttree

import "net/http"

type NotFoundError interface {
	error
	NotFound() bool
}

type notFoundError struct {
}

func (n *notFoundError) Error() string {
	return "not found"
}

func (n *notFoundError) NotFound() bool {
	return true
}

func (n *notFoundError) StatusCode() int {
	return http.StatusNotFound
}

var _ NotFoundError = (*notFoundError)(nil)
var notFoundErr = &notFoundError{}
