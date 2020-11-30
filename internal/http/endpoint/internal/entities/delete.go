package entities

type DeleteRequest interface {
	GetVal() int
}
type DeleteResponse struct {
}

func (r *DeleteResponse) IsDeleteResponse() bool {
	return true
}
