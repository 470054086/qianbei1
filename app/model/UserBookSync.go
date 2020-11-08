package model

import (
	"github.com/Gre-Z/common/jtime"
	errors "github.com/pkg/errors"
	"qianbei.com/core"
)

type UserBookSync struct {
	ID        int            `gorm:"primary_key" json:"id"`
	UserId    int            `gorm:"default:0" json:"-"`
	BookId    int            `gorm:"default:0" json:"book_id"`
	IsDel     int            `gorm:"default:0" json:"-"`
	CreatedAt jtime.JsonTime `json:"created_at"`
}

func (UserBookSync) TableName() string {
	return "user_book_sync"
}

// 根据用户userId和bookId查询数据
func (u *UserBookSync) getFirstByUserIdAndBookId(userId, bookId int) (*UserBookSync, error) {
	var r []*UserBookSync
	if err := core.Db().Where("user_id = ? and book_id = ?", userId, bookId).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err,"")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

// 根据bookId查询多行数据
func (u *UserBookSync) GetUserIdsByBookID(bookId int) ([]*UserBookSync, error) {
	var r []*UserBookSync
	if err := core.Db().Where(" book_id = ?", bookId).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r, nil
	}
	return nil, nil
}
