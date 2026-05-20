package pagination

type Collection[T any] struct {
	Items      []T `json:"items"`
	TotalCount int `json:"total_count"`
}

type PageParams struct {
	Page  int
	Limit int
}

func (p PageParams) Offset() int {
	return (p.Page - 1) * p.Limit
}
