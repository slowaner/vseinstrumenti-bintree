package encoders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/slowaner/vseinstrumenti-bintree/internal/http/server/internal/entities"
)

type appendRequest struct {
	Val int `json:"val"`
}

func DecodeAppendRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	appendRequest := new(appendRequest)
	err = json.NewDecoder(req.Body).Decode(appendRequest)
	if err != nil {
		return
	}

	appendReq := entities.AppendRequest{}
	appendReq.SetVal(appendRequest.Val)

	request = &appendReq
	return
}

func EncodeAppendResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) (err error) {
	_, ok := resp.(entities.AppendResponse)
	if !ok {
		err = errors.New("wrong service response")
		return
	}

	return
}
