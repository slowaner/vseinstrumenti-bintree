package entities

type DeleteRequest struct {
	val int
}

func (r *DeleteRequest) SetVal(val int) {
	r.val = val
}

func (r *DeleteRequest) GetVal() int {
	return r.val
}

type DeleteResponse interface {
}
