package model

import (
	"github.com/Gre-Z/common/jtime"
	"github.com/jinzhu/gorm"
	errors "github.com/pkg/errors"
	"qianbei.com/constat"
	"qianbei.com/core"
)

type UserAccountYm struct {
	ID            int            `gorm:"primary_key" json:"-"`
	BookId        int            `json:"book_id"`
	UserId        int            `gorm:"default:0" json:"-"`
	PayAccount    int            `gorm:"default:0" json:"pay_account"`
	IncomeAccount int            `gorm:"default:0" json:"income_account"`
	IsDel         int            `gorm:"default:0" json:"-"`
	Ym            int            `gorm:"default:0" json:"ym"`
	CreatedAt     jtime.JsonTime `json:"created_at"`
	UpdatedAt     jtime.JsonTime `json:"updated_at"`
}

func (UserAccountYm) TableName() string {
	//实现TableName接口，以达到结构体和表对应，如果不实现该接口，gorm会自动扩展表名为articles（结构体+s）
	return "user_account_ym"
}

// 增加总支出
func (a UserAccountYm) TransactionIncrPay(tx *gorm.DB, userId int, bookId int, account int) (*UserAccountYm, error) {
	r := &UserAccountYm{}
	if err := tx.Where("user_id=? and book_id=?", userId, bookId).First(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	r.PayAccount = r.PayAccount + account
	if err := tx.Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 增加总收入
func (a UserAccountYm) TransactionIncome(tx *gorm.DB, userId int, bookId int, account int) (*UserAccountYm, error) {
	r := &UserAccountYm{}
	if err := tx.Where("user_id=? and book_id=?", userId, bookId).First(r).Error; err != nil {
		return nil, err
	}
	r.IncomeAccount = r.IncomeAccount + account
	if err := tx.Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 根据年月
func (a UserAccountYm) DecrPayYm(userId int, bookId int, account int, ym int) (*UserAccountYm, error) {
	r := &UserAccountYm{}
	if err := core.Db().Where("user_id=? and book_id=? and ym = ? and is_del=?", userId, bookId, ym, constat.NOT_DELETE).
		First(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	r.PayAccount = r.PayAccount - account
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 根据年月
func (a UserAccountYm) DecrIncomeYm(userId int, bookId int, account int, ym int) (*UserAccountYm, error) {
	r := &UserAccountYm{}
	if err := core.Db().Where("user_id=? and book_id=? and ym = ? and is_del = ?", userId, bookId, ym, constat.NOT_DELETE).
		First(r).Error; err != nil {
		return nil, err
	}
	r.IncomeAccount = r.IncomeAccount - account
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

func (a UserAccountYm) DeleteByUserAndBookIdAndYm(userId int, bookId int, ym int) (*UserAccountYm, error) {
	r := &UserAccountYm{}
	if err := core.Db().Where("user_id=? and book_id=? and ym = ? and is_del=?", userId, bookId, ym, constat.NOT_DELETE).
		First(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	r.ID = constat.DELETE
	if err := core.Db().Save(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

func (a *UserAccountYm) DeleteByUserAndBookId(userId int, bookId int) error {
	if err := core.Db().Model(&UserAccountYm{}).Where("user_id=? and book_id=? and is_del=?", userId, bookId, constat.NOT_DELETE).
		Update(&UserAccountYm{IsDel: constat.DELETE}).Error; err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
