package schedule
//
//import (
//	"qianbei.com/app/model"
//	"qianbei.com/constat"
//	"qianbei.com/core"
//	"strconv"
//)
//
//var (
//	tag = "[deletePayRecord]"
//)
//
//type deletePayParams struct {
//	userId int //用户id
//	bookId int //账本id
//}
//
//// 删除用户流水的schedule
//type deletePay struct {
//	deleteChan    chan *deletePayParams
//	payRecord     *model.PayRecord
//	userAccountYm *model.UserAccountYm
//	isBool        bool
//}
//
//var g_SCHEDULE_DEL_PAY *deletePay
//
//func NewDeletePay(nums int) *deletePay {
//	g_SCHEDULE_DEL_PAY = &deletePay{
//		deleteChan:    make(chan *deletePayParams, nums),
//		payRecord:     &model.PayRecord{},
//		userAccountYm: &model.UserAccountYm{},
//	}
//	return g_SCHEDULE_DEL_PAY
//}
//func GetDelPay() *deletePay {
//	return g_SCHEDULE_DEL_PAY
//}
//
//// goruntime 回调函数
//var defDel = func(v *deletePayParams) {
//	r := recover()
//	core.QLog().Error(map[string]interface{}{
//		"goruntime": "添加的goruntime发生错误",
//		"data":      v,
//		"error":     r,
//	})
//}
//
//func (d *deletePay) Run() {
//	if d.isBool == true {
//		return
//	}
//	go func() {
//		for v := range d.deleteChan {
//			defer defDel(v)
//			go func() {
//				d.delete(v)
//			}()
//		}
//	}()
//	d.isBool = true
//}
//
//func (d *deletePay) Schedule(userId, bookId int) {
//	r := &deletePayParams{
//		userId: userId,
//		bookId: bookId,
//	}
//	d.deleteChan <- r
//}
//
//func (d *deletePay) delete(r *deletePayParams) {
//	// 获取账本中的所有数据
//	payRecords, err := d.payRecord.GetListByUserIdAndBookIdAndSync(r.userId, constat.MAIN_BOOK_ID, r.bookId)
//	if err != nil {
//		core.QLog().Error(map[string]interface{}{
//			"tag":  tag,
//			"msg":  err,
//			"data": r,
//		})
//		return
//	}
//	if payRecords == nil {
//		core.QLog().Info(map[string]interface{}{
//			"tag":  tag,
//			"data": r,
//			"msg":  "统计不存在数据",
//		})
//		return
//	}
//
//	// 进行统计 分组进行循环
//	totalPay := 0
//	totalIncome := 0
//	YMtotalPay := map[string]map[string]int{}
//	YMtotalIncome := map[string]map[string]int{}
//	payIds := []int{}
//
//
//	// 先进行年月的修改操作
//	for k, v := range YMtotalPay {
//		ym, err := strconv.Atoi(k)
//		if err != nil {
//			continue
//		}
//		_, err = d.userAccountYm.DecrPayYm(v["user_id"],constat.MAIN_BOOK_ID, v["total"], ym)
//		if err != nil {
//			core.QLog().Error(map[string]interface{}{
//				"tag":  tag,
//				"msg":  "进行年月日更新失败",
//				"err":  err,
//				"data": v,
//			})
//			continue
//		}
//	}
//	// 进行收入的更新
//	for k, v := range YMtotalIncome {
//		ym, err := strconv.Atoi(k)
//		if err != nil {
//			continue
//		}
//		_, err = d.userAccountYm.DecrIncomeYm(v["user_id"], constat.MAIN_BOOK_ID, v["total"], ym)
//		if err != nil {
//			core.QLog().Error(map[string]interface{}{
//				"tag":  tag,
//				"msg":  "进行年月日更新失败",
//				"err":  err,
//				"data": v,
//			})
//			continue
//		}
//	}
//	// 在进行所有的更新
//	err = d.payRecord.DeleteByIdS(payIds)
//	if err != nil {
//		core.QLog().Error(map[string]interface{}{
//			"tag":  tag,
//			"msg":  "进行每条数据的删除失败",
//			"err":  err,
//			"data": payIds,
//		})
//	}
//	core.QLog().Info(map[string]interface{}{
//		"tag":  tag,
//		"msg":  "删除成功",
//		"err":  nil,
//		"data": r,
//	})
//
//}
