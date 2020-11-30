package entities

type AppendRequest interface {
	GetVal() int
}
type AppendResponse struct {
}

func (r *AppendResponse) IsAppendResponse() bool {
	return true
}
