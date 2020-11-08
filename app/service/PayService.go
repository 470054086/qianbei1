package service

import (
	"qianbei.com/app/inter"
	"qianbei.com/app/model"
	"qianbei.com/app/params/request"
	"qianbei.com/app/service/pay"
	"qianbei.com/constat"
)

type PayService struct {
	payModel *model.PayRecord
}

// 根据类型实例化不同的记账方式
func (p *PayService) CreatePay(types int) inter.PayRecord {
	var res inter.PayRecord
	if types == constat.PAY_TYPE_COMMON {
		res = pay.NewPayCommon()
	} else if types == constat.PAY_TYPE_VOICE {
		res = pay.NewPayVoice()
	}
	return res
}

// 查询多个条件
func (p *PayService) Record(userId int, r *request.Record) (*[]model.PayRecord,int, error) {
	return p.payModel.Record(userId, r)
}
