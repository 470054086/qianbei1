package response

import (
	"qianbei.com/app/cache"
	"qianbei.com/app/model"
)

// 查询操作返回值
type Record struct {
	Lists *[]model.PayRecord `json:"lists"`
	Total BaseTotal          `json:"total"`
}

// 增加操作返回值
type AddRecord struct {
	Amount *model.PayRecord        `json:"amount"`
	Total  *model.UserAccount      `json:"total"`
	MTotal *cache.MonthTotalParams `json:"m_total"`
	Book   *model.Book             `json:"book"`
}
