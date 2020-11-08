package constat

// 抛给用户显示的异常
type ShowErrorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s ShowErrorMsg) Error() string {
	return ""
}

// 生成错误的函数
func NewShowMessage(message string) error {
	return ShowErrorMsg{
		Code:    THROW_EXCEP_CODE,
		Message: message,
	}
}

// 抛出给用户的常量定义
const (
	CREATE_BOOK_FULL       = "最多只允许创建5个账单"
	EXISTS_BOOK_TYPE       = "已存在同类型的账单"
	EXISTS_NAME_TYPE       = "存在同名的账单"
	DELETE_BOOK_MAIN_ERROR = "无法删除主账本"
	USER_EXISTS_BOOK_ERROR = "你不存在此账单,无法更新"
	USER_BOOK_OWNEN_ERROR  = "非本人账单,无法更新"
	USER_JOIN_BOOK_EXISTS  = "你此前已加入此账本"
	USER_JOIN_SYNC_ERROR   = "你此前已加入此账本"

	CATEGORY_NOT_EXIST  = "分类不存在"
	CATEGORY_NAME_EXIST = "分类下已存在名称"
)
