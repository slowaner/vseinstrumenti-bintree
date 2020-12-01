package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/endpoint"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/service"
	"github.com/slowaner/vseinstrumenti-bintree/internal/tree/inttree"
	"github.com/stretchr/testify/assert"
)

func getTestInitialTree() (t inttree.Tree) {
	t = inttree.NewTree([]int{6, 22, 54, 8, 32, 7, 63, 21, 5, 75, 3, 2, 1, 0})
	return
}

func init() {
	logger := log.NewNopLogger()
	t := getTestInitialTree()
	t = inttree.NewLoggingIntTree(log.With(logger, "subsystem", "IntTree"), t)
	s, err := service.NewService(logger, t)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	ep := endpoint.NewEndpoints(s)
	r := NewServer(context.Background(), ep)

	handler := NewLoggingMiddleware(logger, r)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		_ = logger.Log(
			"msg", "starting server",
			"addr", srv.Addr,
		)
		if err := srv.ListenAndServe(); err != nil {
			_ = logger.Log("err", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	go func() {
		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		err = srv.Shutdown(ctx)
		if err != nil {
			_ = level.Error(logger).Log(
				"msg", "can't shutdown server",
				"err", err,
			)
			os.Exit(1)
		}

		_ = logger.Log("msg", "shutting down")
	}()
}

func TestFindNodeNotFound(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/find?val=%d", 23))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "not found", string(bts))
}

func TestFindNodeFound(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/find?val=%d", 32))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "{\"data\":32}\n", string(bts))
}

func TestAppendNode(t *testing.T) {
	resp, err := http.Post("http://localhost:8080/append", "application/json", strings.NewReader("{\"val\":23}"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	assert.Len(t, bts, 0)
}

func TestDeleteNodeNotFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/delete?val=%d", -23), nil)
	assert.NoError(t, err)
	c := http.Client{}
	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "404 page not found\n", string(bts))
}

func TestDeleteNodeFound(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/delete?val=%d", 7), nil)
	assert.NoError(t, err)
	c := http.Client{}
	resp, err := c.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	bts, err := ioutil.ReadAll(resp.Body)
	assert.Len(t, bts, 0)
}
