package model

import (
	"github.com/Gre-Z/common/jtime"
	"github.com/pkg/errors"
	"qianbei.com/constat"
	"qianbei.com/core"
)

type CategoryUser struct {
	ID        int            `gorm:"primary_key" json:"id"`
	UserId    int            `gorm:"default:0" json:"-"`
	Pid       int            `gorm:"default:0" json:"pid"`
	Cid       int            `gorm:"default:0" json:"cid"`
	Name      string         `json:"name"`
	Types     int            `gorm:"primary_key" json:"types"`
	Remark    string         `json:"remark"`
	Sorts     int            `gorm:"default:0" json:"-"`
	IsDel     int            `json:"-"`
	CreatedAt jtime.JsonTime `json:"-"`
	List      []*Category    `json:"list"`
}

func (CategoryUser) TableName() string {
	return "category_user"
}

// 获取全部的分类
func (c *CategoryUser) GetList(userId int) ([]*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("user_id=?", userId).Find(&r).Error
		err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

// 获取全部的分类
func (c *CategoryUser) GetFirstById(id int) (*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("id=?", id).Find(&r).Error
		err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func (c *CategoryUser) GetFirstByCId(cid int) (*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("cid=? and is_del = ?", cid, constat.NOT_DELETE).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func (c *CategoryUser) GetFirstByCid(cid int, userId int) (*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("cid=? and user_id = ? ", cid, userId).Find(&r).Error
		err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func (c *CategoryUser) GetFirstByCNameAndUserIdAndCid(name string, cid, userId int) (*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("name=? and cid=? and user_id=? and is_del=? ", name, cid, userId, constat.NOT_DELETE).Find(&r).Error
		err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func (c *CategoryUser) GetFirstByCNameAndUserIdAndPid(name string, userId int, pid int, types int) (*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("name=? and  user_id=? and pid= ? and is_del=? and types =? ", name, userId, pid, constat.NOT_DELETE, types).
		Find(&r).Error
		err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func (c *CategoryUser) GetFirstByCNameAndUserIdAndPidDel(name string, userId int, pid int, types int) (*CategoryUser, error) {
	var r []*CategoryUser
	if err := core.Db().Where("name=? and  user_id=? and pid= ? and is_del=? and types =? ", name, userId, pid, constat.DELETE, types).
		Find(&r).Error
		err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}
