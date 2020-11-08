package service

import (
	"github.com/jinzhu/gorm"
	errors "github.com/pkg/errors"
	"qianbei.com/app/cache"
	"qianbei.com/app/model"
	"qianbei.com/app/params/request"
	"qianbei.com/constat"
	"qianbei.com/core"
	"qianbei.com/util"
	"strconv"
)

type BookService struct {
	bookModel       *model.Book
	bookUserModel   *model.UserBook
	payRecord       *model.PayRecord
	userAccount     *model.UserAccount
	userAccountYm   model.UserAccountYm
	cacheMonthTotal cache.MonthTotal
}

// goruntime 回调函数
var defDel = func(userId int, book *model.Book) {
	r := recover()
	if r != nil {
		core.QLog().Error(map[string]interface{}{
			"goruntime": "添加的goruntime发生错误",
			"data": map[string]interface{}{
				"user_id": userId,
				"book":    book,
			},
			"error": r,
		})
	}

}

const (
	TAG_DEL        = "[DelPayRecord]"
	BOOK_OWNER_MAX = 5 // 用户最多拥有的整本书
)

// 获取用户账本
func (b *BookService) GetList(userId int) ([]*model.Book, error) {
	owner, err := b.bookUserModel.GetUserOwner(userId)
	if err != nil {
		return nil, err
	}
	// 如果为空的话 直接返回
	if len(owner) == 0 {
		return nil, nil
	}
	var idsString = []string{}
	for _, v := range owner {
		idsString = append(idsString, strconv.Itoa(v.ID))
	}
	lists, err := b.bookModel.GetListByIds(idsString)
	if err != nil {
		return nil, err
	}
	return lists, nil
}

// 添加账本
func (b *BookService) Add(userId int, r *request.AddBook) (*model.Book, error) {
	// 判断用户所拥有的账本数量
	owner, err := b.bookUserModel.GetUserOwner(userId)
	if err != nil {
		return nil, err
	}
	// 判断是否存在这个类型的账本
	book, err := b.bookModel.ExistNameBook(userId, r.Name)
	if err != nil {
		return nil, err
	}
	if book {
		return nil, constat.NewShowMessage(constat.EXISTS_NAME_TYPE)
	}

	if len(owner) >= BOOK_OWNER_MAX {
		return nil, constat.NewShowMessage(constat.CREATE_BOOK_FULL)
	}
	var m *model.Book
	// 创建账本 使用事务
	err = core.Db().Transaction(func(tx *gorm.DB) error {
		// 创建账本
		m = &model.Book{
			Name:       r.Name,
			Desc:       r.Desc,
			IsSyncMain: *r.IsSyncMain,
			BgImage:    r.BgImage,
			CreatedId:  userId,
		}
		if err := tx.Create(m).Error; err != nil {
			return errors.Wrap(err, "")
		}
		// 创建账本统计
		t := &model.UserAccount{
			BookId: m.ID,
			UserId: userId,
		}
		if err := tx.Create(t).Error; err != nil {
			return errors.Wrap(err, "")
		}
		// 加入用户
		u := &model.UserBook{
			ID:     0,
			UserId: userId,
			BookId: m.ID,
		}
		if err := tx.Create(u).Error; err != nil {
			return errors.Wrap(err, "")
		}
		// 如果选择了sync 创建sync
		if *r.IsSyncMain == constat.IS_MAIN_SYNC {
			sync := &model.UserBookSync{
				UserId: userId,
				BookId: m.ID,
			}
			if err := tx.Create(sync).Error; err != nil {
				return errors.Wrap(err, "")
			}
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return m, nil
}

// 更新数据
func (b *BookService) Update(userId int, r *request.UpdateBook) (*model.Book, error) {
	userInfo, err := b.bookModel.GetFirstById(r.Id)
	if err != nil {
		return nil, err
	}
	// 判断是否存在
	if userInfo == nil {
		return nil, constat.NewShowMessage(constat.USER_EXISTS_BOOK_ERROR)
	}
	if userInfo.CreatedId != userId {
		return nil, constat.NewShowMessage(constat.USER_BOOK_OWNEN_ERROR)
	}
	// 如果名称进行修改的话
	if userInfo.Name != r.Name {
		existBook, err := b.bookModel.ExistNameBook(userId, r.Name)
		if err != nil {
			return nil, err
		}
		if existBook {
			return nil, constat.NewShowMessage(constat.EXISTS_NAME_TYPE)
		}
	}
	// 修改操作
	if r.Name != "" {
		userInfo.Name = r.Name
	}
	if r.Desc != "" {
		userInfo.Desc = r.Desc
	}
	if r.BgImage != "" {
		userInfo.BgImage = r.BgImage
	}
	// 修改数据
	if err := core.Db().Save(userInfo).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return userInfo, err
}

// 加入用户
func (b *BookService) UserJoin(userId int, r *request.JoinUserRequest) (*model.Book, error) {
	// 判断是否存在账本
	book, err := b.bookModel.GetFirstById(r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	// 判断是否存在
	if book == nil {
		return nil, constat.NewShowMessage(constat.USER_EXISTS_BOOK_ERROR)
	}
	// 判断是否是同步账本
	if book.IsSyncMain == constat.NOT_MAIN_SYNC && *r.IsSyncMain == constat.IS_MAIN_SYNC {
		return nil, constat.NewShowMessage(constat.USER_JOIN_SYNC_ERROR)
	}
	// 判断用户所拥有的账本数量
	owner, err := b.bookUserModel.GetUserOwner(userId)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(owner) >= BOOK_OWNER_MAX {
		return nil, constat.NewShowMessage(constat.CREATE_BOOK_FULL)
	}
	// 判断用户是否加入过此账本
	for _, v := range owner {
		if v.UserId == userId {
			return nil, constat.NewShowMessage(constat.USER_JOIN_BOOK_EXISTS)
		}
	}

	err = core.Db().Transaction(func(tx *gorm.DB) error {
		// 创建账本统计
		t := &model.UserAccount{
			BookId: r.Id,
			UserId: userId,
		}
		if err := tx.Create(t).Error; err != nil {
			return errors.Wrap(err, "")
		}
		// 加入用户
		u := &model.UserBook{
			ID:     0,
			UserId: userId,
			BookId: r.Id,
		}
		// 如果选择了同步的话
		if *r.IsSyncMain == constat.IS_MAIN_SYNC {
			sync := &model.UserBookSync{
				UserId: userId,
				BookId: r.Id,
			}
			if err := tx.Create(sync).Error; err != nil {
				return errors.Wrap(err, "")
			}
		}
		if err := tx.Create(u).Error; err != nil {
			return errors.Wrap(err, "")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	// 返回当前用户的账本
	return b.bookModel.GetFirstById(r.Id)
}

func (b *BookService) DeleteBook(userId int, r *request.DeleteBook) error {

	// 判断是否存在账本
	book, err := b.bookModel.GetFirstById(r.Id)
	if err != nil {
		return err
	}
	// 判断是否存在
	if book == nil {
		return constat.NewShowMessage(constat.USER_EXISTS_BOOK_ERROR)
	}
	if book.IsMain == constat.IS_MAIN_BOOK_ID {
		return constat.NewShowMessage(constat.DELETE_BOOK_MAIN_ERROR)
	}

	// 如果是创建者的话 直接删除
	if book.CreatedId == userId {
		err := core.Db().Transaction(func(tx *gorm.DB) error {
			// 修改book
			book.IsDel = constat.DELETE
			if err := tx.Save(book).Error; err != nil {
				return err
			}
			// 修改userBook
			userBook, err := b.bookUserModel.GetFirstByBookIdAndUserId(book.ID, userId)
			if err != nil {
				return errors.Wrap(err, "")
			}
			if userBook == nil {
				return constat.NewShowMessage(constat.USER_JOIN_BOOK_EXISTS)
			}
			userBook.IsDel = 1
			if err := tx.Save(userBook).Error; err != nil {
				return errors.Wrap(err, "")
			}
			return nil
		})
		if err != nil {
			return errors.Wrap(err, "")
		}
	} else {
		// 如果是成员的话
		userBook, err := b.bookUserModel.GetFirstByBookIdAndUserId(book.ID, userId)
		if err != nil {
			return errors.Wrap(err, "")
		}
		if userBook == nil {
			return constat.NewShowMessage(constat.USER_JOIN_BOOK_EXISTS)
		}
		userBook.IsDel = 1
		if err := core.Db().Save(userBook).Error; err != nil {
			return errors.Wrap(err, "")
		}
	}
	// 删除账本所有相关的数据
	go func() {
		defer defDel(userId, book)
		b.deleteTotal(userId, book)
		// 如果需要同步删除数据 并且账本是同步账本的话
		// 异步删除流水信息
		if *r.DeleteSync == constat.DELETE_MAIN_SYNC && book.IsSyncMain == constat.IS_MAIN_SYNC {
			b.deleteSync(userId, book)
		}
	}()

	return nil
}

// 删除账本相关的统计数据
func (b *BookService) deleteTotal(userId int, book *model.Book) error {
	// 获取所有的交易流水
	payLists, err := b.payRecord.GetListByUserIdAndBookId(userId, book.ID)
	if err != nil {
		e := errors.Wrap(err, "")
		b.log(userId, book, err, "获取所有账本失败", "error")
		return e
	}
	total, m, _ := b.total(payLists)
	// 删除本月支出的缓存
	for k, v := range total {
		if k == util.GetCurrYmShow() {
			// 修改缓存
			b.cacheMonthTotal.DecrPay(userId, book.ID, v["total"])
		}
	}
	// 修改本月收入的缓存
	for k, v := range m {
		if k == util.GetCurrYmShow() {
			// 修改缓存
			b.cacheMonthTotal.DecrIncome(userId, book.ID, v["total"])
		}
	}

	// 1. 删除所有的记录
	err = b.payRecord.DeleteByUserIdAndBookId(userId, book.ID)
	if err != nil {
		e := errors.Wrap(err, "")
		b.log(userId, book, err, "删除账本失败", "error")
		return e
	}
	//  删除当前账本的统计
	err = b.userAccount.DeleteByUserIdAndBookId(userId, book.ID)
	if err != nil {
		e := errors.Wrap(err, "")
		b.log(userId, book, err, "删除账本失败", "error")
		return e
	}
	// 删除每月的统计信息
	err = b.userAccountYm.DeleteByUserAndBookId(userId, book.ID)
	if err != nil {
		e := errors.Wrap(err, "")
		b.log(userId, book, err, "删除账本失败", "error")
		return e
	}
	b.log(userId, book, err, "成功删除所有账本的记录", "info")
	return nil
}

// 删除同步的流水记录
func (b *BookService) deleteSync(userId int, book *model.Book) {
	// 获取账本中的所有数据
	payRecords, err := b.payRecord.GetListByUserIdAndBookIdAndSync(userId, constat.MAIN_BOOK_ID, book.ID)
	if err != nil {
		b.log(userId, book, err, "获取同步流水记录报错", "error")
		return
	}
	if payRecords == nil {
		b.log(userId, book, err, "同步流水记录记录不存在", "error")
		return
	}
	mPayTotal, mIncomeTotal, pids := b.total(payRecords)

	// 先进行年月的支出操作
	for k, v := range mPayTotal {
		// 如果是当月的话 修改缓存
		if k == util.GetCurrYmShow() {
			// 修改缓存
			b.cacheMonthTotal.DecrPay(userId, book.ID, v["total"])
		} else {
			ym, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			_, err = b.userAccountYm.DecrPayYm(v["user_id"], constat.MAIN_BOOK_ID, v["total"], ym)
			if err != nil {
				b.log(userId, book, err, "进行同步流水年月支出修改失败", "error")
				continue
			}
		}

	}
	// 进行年月的输入
	for k, v := range mIncomeTotal {
		if k == util.GetCurrYmShow() {
			// 修改缓存
			b.cacheMonthTotal.DecrIncome(userId, book.ID, v["total"])
		} else {
			ym, err := strconv.Atoi(k)
			if err != nil {
				continue
			}
			_, err = b.userAccountYm.DecrIncomeYm(v["user_id"], constat.MAIN_BOOK_ID, v["total"], ym)
			if err != nil {
				b.log(userId, book, err, "进行同步流水年月输入修改失败", "error")
				continue
			}
		}
	}
	// 删除所有的相关流水
	err = b.payRecord.DeleteByIdS(pids)
	if err != nil {
		b.log(userId, book, err, "进行同步流水删除所有流水失败", "error")
		return
	}
	b.log(userId, book, err, "进行同步流水删除所有成功", "info")
}

// 按照年月统计所有流失
func (p *BookService) total(records *[]model.PayRecord) (map[string]map[string]int, map[string]map[string]int, []int) {
	totalPay := 0
	totalIncome := 0
	YMtotalPay := map[string]map[string]int{}
	YMtotalIncome := map[string]map[string]int{}
	payIds := []int{}
	for _, v := range *records {
		YmShow := util.GetYmShow(v.CreatedAt)
		// 如果是支出的
		if v.Types == constat.PAY_EXPEND {
			totalPay += v.Account
			// 进行按年月统计
			if _, ok := YMtotalPay[YmShow]; !ok {
				YMtotalPay[YmShow] = map[string]int{
					"total":   v.Account,
					"user_id": v.UserId,
					"book_id": v.BookId,
				}
			} else {
				YMtotalPay[YmShow]["total"] += v.Account
			}
		} else if v.Types == constat.PAY_INCOME {
			totalIncome += v.Account
			if _, ok := YMtotalIncome[YmShow]; !ok {
				YMtotalIncome[YmShow] = map[string]int{
					"total":   v.Account,
					"user_id": v.UserId,
					"book_id": v.BookId,
				}
			} else {
				YMtotalIncome[YmShow]["total"] += v.Account
			}
		}
		payIds = append(payIds, v.ID)
	}
	return YMtotalPay, YMtotalIncome, payIds
}

// 添加日志
func (p *BookService) log(userId int, r *model.Book, err error, msg string, level string) {
	if level == "info" {
		core.QLog().Info(map[string]interface{}{
			"tag": TAG_DEL,
			"data": map[string]interface{}{
				"userId": userId,
				"req":    r,
			},
			"err": nil,
			"msg": msg,
		})
	} else {
		core.QLog().Error(map[string]interface{}{
			"tag": TAG_DEL,
			"data": map[string]interface{}{
				"userId": userId,
				"req":    r,
			},
			"err": err,
			"msg": msg,
		})
	}
}
