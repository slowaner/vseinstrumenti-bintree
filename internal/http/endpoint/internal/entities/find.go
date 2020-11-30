package entities

type FindRequest interface {
	GetVal() int
}

type FindResponse struct {
	data int
}

func (r *FindResponse) SetFoundData(data int) {
	r.data = data
}

func (r *FindResponse) GetFoundData() (data int) {
	return r.data
}
