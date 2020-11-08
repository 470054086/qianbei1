package model

import (
	"github.com/Gre-Z/common/jtime"
	errors "github.com/pkg/errors"
	"qianbei.com/core"
)

type Category struct {
	ID        int            `gorm:"primary_key" json:"id"`
	Pid       int            `gorm:"default:0" json:"pid"`
	Name      string         `json:"name"`
	Types     int            `gorm:"primary_key" json:"types"`
	Remark    string         `json:"remark"`
	Sorts     int            `gorm:"default:0" json:"sorts"`
	Isdel     int            `gorm:"default:0" json:"is_del"` //数据库中不含有此字段 只是为了显示
	CreatedAt jtime.JsonTime `json:"-"`
	List      []*Category    `json:"list"`
}

func (Category) TableName() string {
	return "category"
}

// 排序
type CategorySlice []*Category

func (a CategorySlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a CategorySlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a CategorySlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Sorts < a[i].Sorts
}

// 获取全部的分类
func (c *Category) GetList() ([]*Category, error) {
	var r []*Category
	if err := core.Db().Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return r, nil
}

func (c *Category) GetFirstById(id int) (*Category, error) {
	var r []*Category
	if err := core.Db().Where("id=?", id).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}



func (c *Category) GetFirstByName(name string) (*Category, error) {
	var r []*Category
	if err := core.Db().Where("name=?", name).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}

func (c *Category) GetFirstByNameAndPid(name string, pid int, types int) (*Category, error) {
	var r []*Category
	if err := core.Db().Where("name=? and pid=? and types = ?", name, pid, types).Find(&r).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	if len(r) > 0 {
		return r[0], nil
	}
	return nil, nil
}
