package schedule

import (
	"qianbei.com/app/cache"
	"qianbei.com/app/model"
	"qianbei.com/app/params/request"
	"qianbei.com/core"
)

var (
	tagAdd = "[addPayRecord]"
)

type addPayParams struct {
	userId int
	req    *request.PayRecord
}

// 删除用户流水的schedule
type addPay struct {
	addChan           chan *addPayParams
	payRecord         *model.PayRecord
	userBookSyncModel *model.UserBookSync
	cacheMontotal     *cache.MonthTotal
	userAccount       *model.UserAccount
	isBool            bool
}

var g_SCHEDU_ADD_PAY *addPay

func NewAddPay(nums int) *addPay {
	g_SCHEDU_ADD_PAY = &addPay{
		addChan:           make(chan *addPayParams, nums),
		payRecord:         &model.PayRecord{},
		userBookSyncModel: &model.UserBookSync{},
		cacheMontotal:     &cache.MonthTotal{},
		userAccount:       &model.UserAccount{},
		isBool:            false,
	}
	return g_SCHEDU_ADD_PAY
}
func GetAddPay() *addPay {
	return g_SCHEDU_ADD_PAY
}

var defAdd = func(v *addPayParams) {
	r := recover()
	core.QLog().Error(map[string]interface{}{
		"goruntime": "添加的goruntime发生错误",
		"data":      v,
		"error":     r,
	})
}

func (p *addPay) Run() {
	if p.isBool == true {
		return
	}
	go func() {

	}()
	p.isBool = true
}

func (p *addPay) Schedule(userId int, r *request.PayRecord) {

}


