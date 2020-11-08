package model

import (
	"github.com/Gre-Z/common/jtime"
	errors "github.com/pkg/errors"
	"qianbei.com/constat"
	"qianbei.com/core"
)

type UserBook struct {
	ID        int            `gorm:"primary_key" json:"id"`
	UserId    int            `gorm:"default:0" json:"-"`
	BookId    int            `gorm:"default:0" json:"book_id"`
	IsDel     int            `gorm:"default:0"`
	CreatedAt jtime.JsonTime `json:"created_at"`
	DeletedAt jtime.JsonTime `json:"deleted_at"`
}

func (UserBook) TableName() string {
	return "user_book"
}

func (u *UserBook) GetFirstByBookIdAndUserId(bookId int, userId int) (*UserBook, error) {
	var r []*UserBook
	if err := core.Db().Model(&UserBook{}).
		Where("book_id = ? and user_id = ? and is_del = ?", bookId, userId, constat.NOT_DELETE).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

// 获取用户拥有的账单数量
func (u *UserBook) GetUserOwner(userId int) ([]*UserBook, error) {
	var r []*UserBook
	if err := core.Db().Model(&UserBook{}).
		Where("user_id = ? and is_del = ?", userId, constat.NOT_DELETE).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}
