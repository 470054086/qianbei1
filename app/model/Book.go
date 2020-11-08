package model

import (
	"github.com/Gre-Z/common/jtime"
	errors "github.com/pkg/errors"
	"qianbei.com/constat"
	"qianbei.com/core"
)

type Book struct {
	ID         int            `gorm:"primary_key" json:"id"`
	Name       string         `gorm:"default:''"json:"name"`
	Desc       string         `gorm:"default:''" json:"desc"`
	BgImage    string         `gorm:"default:''"json:"bg_image"`
	IsMain     int            `gorm:"default:''" json:"is_main"`
	IsSyncMain int            `gorm:"default:''" json:"is_sync_main"`
	IsDel      int            `gorm:"default:0" json:"-"`
	CreatedId  int            `json:"created_id"`
	CreatedAt  jtime.JsonTime `json:"created_at"`
	UpdatedAt  jtime.JsonTime `json:"updated_at"`
	DeletedAt  jtime.JsonTime `json:"deleted_at"`
}

func (Book) TableName() string {
	return "book"
}

/**
修改数据
*/
func (b *Book) Update(r *Book) (*Book, error) {
	if err := core.Db().Where("id=? and is_del", r.ID, constat.NOT_DELETE).Update(r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

/**
多个id获取全部的数据
*/
func (b *Book) GetListByIds(ids []string) ([]*Book, error) {
	var r []*Book
	if err := core.Db().Where("id in (?) and is_del=?", ids, constat.NOT_DELETE).
		Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 根据id获取list
func (b *Book) GetFirstById(id int) (*Book, error) {
	var r []*Book
	if err := core.Db().Where("id =? and is_del=?", id, constat.NOT_DELETE).
		Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

// 根据userId获取list
func (b *Book) GetListByUserId(userId int) ([]*Book, error) {
	var r []*Book
	if err := core.Db().Where("created_id =? and is_del=?", userId, constat.NOT_DELETE).
		Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 根据userId和Type判断是否存在
func (b *Book) ExistTypeBook(userId int, bookType int) (bool, error) {
	var r []*Book
	if err := core.Db().Model(&Book{}).
		Where("created_id=? and book_type=? and is_del=?", userId, bookType, constat.NOT_DELETE).Find(&r).Error; err != nil {
		return false, errors.Wrap(err, "")
	}
	if len(r) != 0 {
		return true, nil
	}
	return false, nil
}

// 判断name是否存在
func (b *Book) ExistNameBook(userId int, name string) (bool, error) {
	var r []*Book
	if err := core.Db().Model(&Book{}).
		Where("created_id=? and name=? and is_del=?", userId, name, constat.NOT_DELETE).Find(&r).Error; err != nil {
		return false, errors.Wrap(err, "")
	}
	if len(r) != 0 {
		return true, nil
	}
	return false, nil
}
