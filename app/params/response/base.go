package response

/**
	 分页返回的基本类型
 */
type BaseTotal struct {
	Page int `json:"page"`
	PageSize int `json:"page_size"`
	Count int `json:"count"`
}
