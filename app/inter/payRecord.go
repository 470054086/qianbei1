package inter

import (
	"qianbei.com/app/params/request"
	"qianbei.com/app/params/response"
)

// 流水写入表
type PayRecord interface {
	AddRecord(userId int, r *request.PayRecord) (*response.AddRecord, error) // 添加流水
	UpdateRecord()                                                              //修改流水
	DeleteRecord()                                                              //删除流水
}
