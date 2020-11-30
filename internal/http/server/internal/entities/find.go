package entities

type FindRequest struct {
	val int
}

func (r *FindRequest) SetVal(val int) {
	r.val = val
}

func (r *FindRequest) GetVal() int {
	return r.val
}

type FindResponse interface {
	GetFoundData() (data int)
}
