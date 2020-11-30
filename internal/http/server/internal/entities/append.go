package entities

type AppendRequest struct {
	val int
}

func (r *AppendRequest) SetVal(val int) {
	r.val = val
}

func (r *AppendRequest) GetVal() int {
	return r.val
}

type AppendResponse interface {
}
