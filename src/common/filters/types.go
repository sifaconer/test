package filters

type FilterParams struct {
	Filters    string
	Sort       string
	Pagination *PaginationParams
}

type PaginationParams struct {
	Page int
	Size int
}
type PaginationInfo struct {
	Page      int `json:"page"`
	Size      int `json:"size"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page,omitempty"`
}

type FilterCondition struct {
	Field    string
	Operator string
	Value    string
}

type FilterGroup struct {
	Type       string
	Conditions []FilterCondition
	Groups     []FilterGroup
}

type SortField struct {
	Field string
	Order string // "ASC" o "DESC"
}
