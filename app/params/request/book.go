package request

// 创建账本
type AddBook struct {
	Name       string `json:"name" binding:"required"`
	Desc       string `json:"desc"`
	IsSyncMain *int   `json:"is_sync_main" binding:"required"` //是否同步主账单
	BgImage    string `json:"bg_image"`
}

// 修改账本
type UpdateBook struct {
	IdRequest
	Name    string `json:"name"`
	Desc    string `json:"desc"`
	BgImage string `json:"bg_image"`
}

// 加入账本
type JoinUserRequest struct {
	IdRequest
	IsSyncMain *int `json:"is_sync_main" binding:"required"` //是否同步主账单
	UserId  int `json:"user_id"`
}

// 删除账本
type DeleteBook struct {
	IdRequest
	DeleteSync *int `json:"delete_sync" binding:"required"` //是否删除到主账单的数据
	UserId  int `json:"user_id"`
}
