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
	Order  string `json:"order,omitempty,omitzero"`
	Search string `json:"search,omitempty,omitzero"`
	Page   int    `json:"page,omitempty"`
	Size   int    `json:"size,omitempty"`
}

func (qp *QueryParams) String() string {
	return ""
}

type APIError struct {
	Field   string `json:"field,omitempty,omitzero"`
	Message string `json:"message"`
}
