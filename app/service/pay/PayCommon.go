package pay

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"qianbei.com/app/cache"
	"qianbei.com/app/model"
	"qianbei.com/app/params/request"
	"qianbei.com/app/params/response"
	"qianbei.com/constat"
	"qianbei.com/core"
)

var (
	tagAdd   = "[addPayRecord]"
	tagAddGo = "[addPayRecordGoruntine]"
	// 定义捕捉错误
	defAdd = func(userId int, v *request.PayRecord) {
		r := recover()
		core.QLog().Error(map[string]interface{}{
			"msg":   tagAddGo,
			"user":  userId,
			"data":  v,
			"error": r,
		})
	}
)

type payCommon struct {
	payModel     *model.PayRecord
	userAccount  *model.UserAccount
	bookModel    *model.Book
	mUserAccount cache.MonthTotal
	userBookSync *model.UserBookSync
}

func NewPayCommon() *payCommon {
	return &payCommon{}
}

// 添加操作
func (p *payCommon) AddRecord(userId int, r *request.PayRecord) (*response.AddRecord, error) {
	return p.addPay(userId, r)
}

func (p *payCommon) UpdateRecord() {
}

func (p *payCommon) DeleteRecord() {
}

func (p *payCommon) addPay(userId int, r *request.PayRecord) (*response.AddRecord, error) {
	record, err := p.AddInsertRecord(userId, r, 0)
	if err != nil {
		return nil, err
	}
	// 如果是同步账本 并且不是主账本
	if record.Book.IsSyncMain == constat.DELETE_MAIN_SYNC && record.Book.ID != constat.MAIN_BOOK_ID {
		go func(userId int, r *request.PayRecord) {
			defAdd(userId, r)
			p.addSyncRecord(r)
		}(userId, r)
	}
	return record, nil
}

// 添加同步账本
func (p *payCommon) addSyncRecord(r *request.PayRecord) {
	// 查询当前账本下面是否存在用户
	userSyncBook, err := p.userBookSync.GetUserIdsByBookID(r.BookId)
	if err != nil {
		core.QLog().Error(map[string]interface{}{
			"tag":  tagAdd,
			"data": r,
			"msg":  "查询是否有同步用户存在报错",
		})
		return
	}
	if userSyncBook == nil {
		core.QLog().Info(map[string]interface{}{
			"tag":  tagAdd,
			"data": r,
			"msg":  "不存在同步用户",
		})
		return
	}
	flag := true
	// 同步的时候 许改bookID为主账单
	bid := r.BookId
	r.BookId = constat.MAIN_BOOK_ID
	for _, v := range userSyncBook {
		// 调取添加流水的函数
		_, err := p.AddInsertRecord(v.UserId, r, bid)
		if err != nil {
			flag = false
			continue
		}
	}
	if flag {
		core.QLog().Info(map[string]interface{}{
			"tag":  tagAdd,
			"data": r,
			"msg":  "当前流水同步成功",
		})
	}
}

// 添加流水记录
func (p *payCommon) AddInsertRecord(userId int, r *request.PayRecord, syncBookId int) (*response.AddRecord, error) {
	// 定义变量
	var m *model.PayRecord
	var income *model.UserAccount
	var mTotal *cache.MonthTotalParams
	var book *model.Book
	// 判断是否存在
	book, err := p.bookModel.GetFirstById(r.BookId)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if book == nil {
		return nil, errors.New("查询不到账本信息")
	}
	// 如果是主账本的话 同步账本id为0 非主账本 同步账本id从数据库获取
	err = core.Db().Transaction(func(tx *gorm.DB) error {
		// 1. 添加流水记录
		m = &model.PayRecord{
			UserId:     userId,
			Category:   r.Category,
			BookId:     r.BookId,
			Account:    r.Account,
			PlatForm:   r.Platform,
			Types:      r.Types,
			SyncBookId: syncBookId,
			UploadType: r.UploadType,
			Memo:       r.Memo,
		}
		err := tx.Create(m).Error
		if err != nil {
			e := errors.Wrap(err, "")
			p.log(userId, r, e, "添加流水记录出错", "error")
			return errors.Wrap(err, "")
		}
		// 修改总支出函数
		if r.Types == constat.PAY_EXPEND {
			income, err = p.userAccount.TransactionIncrPay(tx, userId, r.BookId, r.Account)
			if err != nil {
				e := errors.Wrap(err, "")
				p.log(userId, r, e, "修改总支出出错", "error")
				return e
			}
		} else {
			income, err = p.userAccount.TransactionIncome(tx, userId, r.BookId, r.Account)
			if err != nil {
				e := errors.Wrap(err, "")
				p.log(userId, r, e, "修改总收入出错", "error")
				return e
			}
		}
		// 增加本月支出
		mTotal, err = p.mUserAccount.AddTotal(userId, r.Account, r.BookId, r.Types)
		if err != nil {
			e := errors.Wrap(err, "")
			p.log(userId, r, e, "增加本月收入/支出出错", "error")
			return e
		}
		return nil
	})
	if err != nil {
		e := errors.Wrap(err, "添加数据发生错误")
		p.log(userId, r, e, "数据库事务报错", "error")
		return nil, e

	}
	p.log(userId, r, nil, "添加数据成功", "info")
	return &response.AddRecord{
		Amount: m,
		Total:  income,
		MTotal: mTotal,
		Book:   book,
	}, nil
}

// 添加日志
func (p *payCommon) log(userId int, r *request.PayRecord, err error, msg string, level string) {
	if level == "info" {
		core.QLog().Info(map[string]interface{}{
			"tag": tagAdd,
			"data": map[string]interface{}{
				"userId": userId,
				"req":    r,
			},
			"err": nil,
			"msg": msg,
		})
	} else {
		core.QLog().Error(map[string]interface{}{
			"tag": tagAdd,
			"data": map[string]interface{}{
				"userId": userId,
				"req":    r,
			},
			"err": err,
			"msg": msg,
		})
	}
}
