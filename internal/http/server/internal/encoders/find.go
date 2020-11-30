package encoders

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/slowaner/vseinstrumenti-bintree/internal/http/server/internal/entities"
)

type findResp struct {
	Data int `json:"data"`
}

func DecodeFindRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
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

	findReq := entities.FindRequest{}
	findReq.SetVal(int(val))

	request = &findReq
	return
}

func EncodeFindResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) (err error) {
	res, ok := resp.(entities.FindResponse)
	if !ok {
		err = errors.New("wrong service response")
		return
	}

	resDt := findResp{
		Data: res.GetFoundData(),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resDt)
	if err != nil {
		return
	}

	return
}
