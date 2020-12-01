package inttree

import (
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func newTree() Tree {
	logger := log.NewNopLogger()
	t := NewTree([]int{6, 22, 54, 8, 32, 7, 63, 21, 5, 75, 3, 2, 1, 0})
	t = NewLoggingIntTree(logger, t)
	return t
}

func TestNotFoundError(t *testing.T) {
	type httpError interface {
		error
		StatusCode() int
	}
	notFoundErrorType := new(NotFoundError)
	httpErrorType := new(httpError)
	assert.Implements(t, notFoundErrorType, notFoundErr)
	assert.Implements(t, httpErrorType, notFoundErr)
	assert.EqualError(t, notFoundErr, "not found")
	assert.Equal(t, true, notFoundErr.NotFound())
	assert.Equal(t, http.StatusNotFound, notFoundErr.StatusCode())
}

func TestFindNodeNotFound(t *testing.T) {
	tr := newTree()
	foundData, err := tr.Find(23)
	assert.EqualError(t, err, notFoundErr.Error())
	assert.Equal(t, 0, foundData)
}

func TestFindNodeFound(t *testing.T) {
	tr := newTree()
	foundData, err := tr.Find(32)
	assert.NoError(t, err)
	assert.Equal(t, 32, foundData)
}

func TestAppendNode(t *testing.T) {
	tr := newTree()
	err := tr.Append(23)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteNodeNotFound(t *testing.T) {
	tr := newTree()
	err := tr.Delete(23)
	assert.EqualError(t, err, notFoundErr.Error())
}

func TestDeleteNodeFound(t *testing.T) {
	tr := newTree()
	err := tr.Delete(7)
	assert.NoError(t, err)
}

func TestLNodeRNodeDeleting(t *testing.T) {
	tr := newTree()
	err := tr.Delete(6)
	assert.NoError(t, err)
	err = tr.Delete(5)
	assert.NoError(t, err)
	err = tr.Delete(22)
	assert.NoError(t, err)
	err = tr.Delete(3)
	assert.NoError(t, err)
	err = tr.Delete(75)
	assert.NoError(t, err)
}
