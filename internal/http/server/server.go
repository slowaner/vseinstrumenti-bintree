package server

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/server/internal/encoders"
)

// endpoints contains endpoints for service
type endpoints interface {
	// GetFindEndpoint returns endpoint for find tree node
	GetFindEndpoint() endpoint.Endpoint
	// GetAppendEndpoint returns endpoint for append tree node
	GetAppendEndpoint() endpoint.Endpoint
	// GetDeleteEndpoint returns endpoint for delete tree node
	GetDeleteEndpoint() endpoint.Endpoint
}

func NewServer(ctx context.Context, endpoints endpoints) *mux.Router {
	r := mux.NewRouter()

	findHandler := ht.NewServer(
		endpoints.GetFindEndpoint(),
		encoders.DecodeFindRequest,
		encoders.EncodeFindResponse,
	)
	appendHandler := ht.NewServer(
		endpoints.GetAppendEndpoint(),
		encoders.DecodeAppendRequest,
		encoders.EncodeAppendResponse,
	)
	deleteHandler := ht.NewServer(
		endpoints.GetDeleteEndpoint(),
		encoders.DecodeDeleteRequest,
		encoders.EncodeDeleteResponse,
	)

	r.Handle("/find", findHandler).Methods(http.MethodGet).Queries("val", "{val:\\d+}")
	r.Handle("/append", appendHandler).Methods(http.MethodPost).Headers("Content-Type", "application/json")
	r.Handle("/delete", deleteHandler).Methods(http.MethodDelete).Queries("val", "{val:\\d+}")

	return r
}
