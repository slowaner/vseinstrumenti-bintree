package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/endpoint"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/server"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/service"
	"github.com/slowaner/vseinstrumenti-bintree/internal/tree/inttree"
)

func getInitialTree() (t inttree.Tree, err error) {
	bts, err := ioutil.ReadFile("initial-tree.json")
	if err != nil {
		return
	}

	var elems []int
	err = json.Unmarshal(bts, &elems)
	if err != nil {
		return
	}

	t = inttree.NewTree(elems)
	return
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	logger := log.NewJSONLogger(os.Stdout)
	t, err := getInitialTree()
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	t = inttree.NewLoggingIntTree(log.With(logger, "subsystem", "IntTree"), t)
	s, err := service.NewService(logger, t)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	ep := endpoint.NewEndpoints(s)
	r := server.NewServer(context.Background(), ep)

	handler := server.NewLoggingMiddleware(logger, r)

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

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
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
	os.Exit(0)
}
