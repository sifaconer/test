package common

type Response[T any] struct {
	Status     string       `json:"status"`
	Code       int          `json:"code"`
	Message    string       `json:"message"`
	Data       T            `json:"data,omitempty,omitzero"`
	Pagination *Pagination  `json:"pagination,omitempty,omitzero"`
	Query      *QueryParams `json:"query"`
	Errors     []APIError   `json:"errors,omitempty,omitzero"`
}

type Pagination struct {
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}

type QueryParams struct {
	Field  string `json:"field,omitempty,omitzero"`
	Filter string `json:"filter,omitempty,omitzero"`
	Sort  string `json:"sort,omitempty,omitzero"`
	Page   int    `json:"page,omitempty" validate:"min=1"`
	Size   int    `json:"size,omitempty"`
}

func (qp *QueryParams) Default() {
	if qp == nil {
		return
	}
	if qp.Page == 0 {
		qp.Page = 1
	}
	if qp.Size == 0 {
		qp.Size = 10
	}
}

func (qp *QueryParams) IsEmpty() bool {
	return qp == nil || qp.Field == "" && qp.Filter == "" && qp.Sort == "" && qp.Page == 0 && qp.Size == 0
}

type APIError struct {
	Field   string `json:"field,omitempty,omitzero"`
	Message string `json:"message"`
}
