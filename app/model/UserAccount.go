package model

import (
	"github.com/Gre-Z/common/jtime"
	"github.com/jinzhu/gorm"
	errors "github.com/pkg/errors"
	"qianbei.com/constat"
	"qianbei.com/core"
)

type UserAccount struct {
	ID            int            `gorm:"primary_key" json:"-"`
	BookId        int            `json:"book_id"`
	UserId        int            `gorm:"default:0" json:"-"`
	PayAccount    int            `gorm:"default:0" json:"pay_account"`
	IncomeAccount int            `gorm:"default:0" json:"income_account"`
	IsDel         int            `gorm:"default:0" json:"-"`
	CreatedAt     jtime.JsonTime `json:"created_at"`
	UpdatedAt     jtime.JsonTime `json:"updated_at"`
}

func (UserAccount) TableName() string {
	//实现TableName接口，以达到结构体和表对应，如果不实现该接口，gorm会自动扩展表名为articles（结构体+s）
	return "user_account"
}

// 增加总支出 用于事务
func (a *UserAccount) TransactionIncrPay(tx *gorm.DB, userId int, bookId int, account int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := tx.Where("user_id=? and book_id=? and is_del =? ", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	r.PayAccount = r.PayAccount + account
	if err := tx.Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 增加总收入 用于事务
func (a *UserAccount) TransactionIncome(tx *gorm.DB, userId int, bookId int, account int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := tx.Where("user_id=? and book_id=? and is_del =? ", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, err
	}
	r.IncomeAccount = r.IncomeAccount + account
	if err := tx.Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 增加总支出
func (a *UserAccount) IncrPay(userId int, bookId int, account int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del =? ", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	r.PayAccount = r.PayAccount + account
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 增加总收入
func (a *UserAccount) Income(userId int, bookId int, account int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del = ?", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, err
	}
	r.IncomeAccount = r.IncomeAccount + account
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 增加总支出
func (a *UserAccount) DecrPay(userId int, bookId int, account int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del =? ", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	r.PayAccount = r.PayAccount - account
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 增加总收入
func (a *UserAccount) DecrIncome(userId int, bookId int, account int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del = ?", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, err
	}
	r.IncomeAccount = r.IncomeAccount - account
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

func (a *UserAccount) DecrPayIncome(userId int, bookId int, pay int, income int) (*UserAccount, error) {
	r := &UserAccount{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del = ?", userId, bookId, constat.NOT_DELETE).First(r).Error; err != nil {
		return nil, err
	}
	r.IncomeAccount = r.IncomeAccount - pay
	r.PayAccount = r.PayAccount - income
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

func (a *UserAccount) DeleteByUserIdAndBookId(userId int, bookId int) error {
	if err := core.Db().Model(&UserAccount{}).Where("user_id=? and book_id=? and is_del = ?", userId, bookId, constat.NOT_DELETE).
		Update(&UserAccount{IsDel: constat.DELETE}).Error; err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
