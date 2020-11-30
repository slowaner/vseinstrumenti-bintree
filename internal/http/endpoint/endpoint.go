package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"

	"github.com/slowaner/vseinstrumenti-bintree/internal/http/endpoint/internal/entities"
)

type Service interface {
	Find(ctx context.Context, value int) (val int, err error)
	Append(ctx context.Context, value int) (err error)
	Delete(ctx context.Context, value int) (err error)
}

// Endpoints contains endpoints for service
type Endpoints interface {
	// GetFindEndpoint returns endpoint for find tree node
	GetFindEndpoint() endpoint.Endpoint
	// GetAppendEndpoint returns endpoint for append tree node
	GetAppendEndpoint() endpoint.Endpoint
	// GetDeleteEndpoint returns endpoint for delete tree node
	GetDeleteEndpoint() endpoint.Endpoint
}

var _ Endpoints = (*endpoints)(nil)

type endpoints struct {
	find   endpoint.Endpoint
	append endpoint.Endpoint
	delete endpoint.Endpoint
}

func (e *endpoints) GetFindEndpoint() endpoint.Endpoint {
	return e.find
}

func (e *endpoints) GetAppendEndpoint() endpoint.Endpoint {
	return e.append
}

func (e *endpoints) GetDeleteEndpoint() endpoint.Endpoint {
	return e.delete
}

func makeFindEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		req, ok := request.(entities.FindRequest)
		if !ok {
			err = errors.New("bad request")
			return
		}

		foundData, err := s.Find(ctx, req.GetVal())
		if err != nil {
			return
		}

		findResp := entities.FindResponse{}
		findResp.SetFoundData(foundData)

		resp = &findResp

		return
	}
}

func makeAppendEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		req, ok := request.(entities.AppendRequest)
		if !ok {
			err = errors.New("bad request")
		}

		err = s.Append(ctx, req.GetVal())
		if err != nil {
			return
		}

		resp = &entities.AppendResponse{}

		return
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		req, ok := request.(entities.DeleteRequest)
		if !ok {
			err = errors.New("bad request")
		}

		err = s.Delete(ctx, req.GetVal())
		if err != nil {
			return
		}

		resp = &entities.DeleteResponse{}
		return
	}
}

func NewEndpoints(s Service) Endpoints {
	return &endpoints{
		find:   makeFindEndpoint(s),
		append: makeAppendEndpoint(s),
		delete: makeDeleteEndpoint(s),
	}
}
