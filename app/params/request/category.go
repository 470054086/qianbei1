package request

type CreateCategory struct {
	Name  string `json:"name" binding:"required"`
	Pid   *int    `json:"pid"`
	Sorts *int    `json:"sorts"`
	Types int    `json:"types"`
}

type UpdateCategory struct {
	IdRequest
	Name  string `json:"name"`
	Sorts *int   `json:"sorts"`
}

type DelCategory struct {
	IdRequest
}
