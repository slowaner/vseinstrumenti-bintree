package encoders

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/server/internal/entities"
)

func DecodeDeleteRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	vars := mux.Vars(req)
	valStr, ok := vars["val"]
	if !ok {
		err = errors.New("no value specified")
		return
	}

	val, err := strconv.ParseInt(valStr, 10, 32)
	if err != nil {
		return
	}

	deleteReq := entities.DeleteRequest{}
	deleteReq.SetVal(int(val))

	request = &deleteReq
	return
}

func EncodeDeleteResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) (err error) {
	_, ok := resp.(entities.DeleteResponse)
	if !ok {
		err = errors.New("wrong service response")
		return
	}

	return
}
