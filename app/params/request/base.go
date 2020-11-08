package request

// 分页的请求参数
type BaseRequest struct {
	Page     int `json:"page" binding:"required"`
	PageSize int `json:"page_size" binding:"required"`
}

// 常用id的请求
type IdRequest struct {
	Id int `json:"id" binding:"required"  `
}
