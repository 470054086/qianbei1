package model

import (
	"fmt"
	"github.com/Gre-Z/common/jtime"
	errors "github.com/pkg/errors"
	"qianbei.com/app/params/request"
	"qianbei.com/constat"
	"qianbei.com/core"
)

// 支出数据类型
type PayRecord struct {
	ID         int            `gorm:"primary_key" json:"id"`
	UserId     int            `gorm:"default:0" json:"-"`
	BookId     int            `gorm:"default:0" json:"book_id"`
	Category   int            `gorm:"default:1" json:"category"`
	Account    int            `gorm:"default:0" json:"account"`
	PlatForm   int            `gorm:"default:1" json:"plat_form"`
	Types      int            `gorm:"default:1" json:"types"`
	UploadType int            `gorm:"default:1" json:"upload_type"`
	SyncBookId int            `gorm:"default:0" json:"sync_book_id"` // 同步账本id
	AssetsId   int            `gorm:"default:0" json:"assets_id"`    // 资产id
	Memo       string         `gorm:"default:''" json:"memo"`
	IsDel      int            `gorm:"default:0" json:"-"`
	CreatedAt  jtime.JsonTime `json:"created_at"`
	UpdatedAt  jtime.JsonTime `json:"updated_at"`
}

func (PayRecord) TableName() string {
	//实现TableName接口，以达到结构体和表对应，如果不实现该接口，gorm会自动扩展表名为articles（结构体+s）
	return "pay_record"
}

// 多条件查询
func (p *PayRecord) Record(userId int, r *request.Record) (*[]PayRecord, int, error) {
	res := &[]PayRecord{}
	var counts int
	var db = core.Db().Where("user_id=? and book_id=? and is_del=?", userId, r.BookId, constat.NOT_DELETE)
	// 金額
	if r.StartAccount != 0 && r.EndAccount != 0 {
		db = db.Where("account>=? and account<=? ", r.StartAccount, r.EndAccount)
	}
	// 类型
	if r.Types != 0 {
		db = db.Where("types=?", r.Types)
	}
	// 注释
	if r.Memo != "" {
		db = db.Where("memo like ?", "%"+r.Memo+"%")
	}
	if r.StartTime != "" && r.EndTime != "" {
		db = db.Where("created_at >=? and created_at <=?", r.StartTime, r.EndTime)
	}
	// 排序字段
	if r.Order != "" && r.Desc != "" {
		db = db.Order(fmt.Sprintf("%s %s", r.Order, r.Desc))
	} else {
		db = db.Order("created_at desc")
	}

	if r.Order == "account" {

	}
	if err := db.Offset((r.Page - 1) * r.PageSize).Limit(r.PageSize).Find(res).Error; err != nil {
		return nil, 0, errors.Wrap(err, "")
	}
	if res != nil {
		db.Model(&PayRecord{}).Count(&counts)
	}
	return res, counts, nil
}

func (p *PayRecord) GetListByUserIdAndBookIdAndSync(userId, bookId int, syncBookId int) (*[]PayRecord, error) {
	res := &[]PayRecord{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del = ? and sync_book_id = ?", userId, bookId, constat.NOT_DELETE, syncBookId).
		Find(res).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return res, nil
}

func (p *PayRecord) GetListByUserIdAndBookId(userId, bookId int) (*[]PayRecord, error) {
	res := &[]PayRecord{}
	if err := core.Db().Where("user_id=? and book_id=? and is_del = ?", userId, bookId, constat.NOT_DELETE).
		Find(res).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	return res, nil
}

// 根据ids删除多条记录
func (p *PayRecord) DeleteByIdS(ids []int) error {
	return core.Db().Model(&PayRecord{}).Where("id in (?)", ids).Update(&PayRecord{IsDel: constat.DELETE}).Error
}

func (p *PayRecord) DeleteByUserIdAndBookId(userId, bookId int) error {
	if err := core.Db().Model(&PayRecord{}).Where("user_id =? and book_id = ? and is_del = ?", userId, bookId, constat.NOT_DELETE).
		Update(&PayRecord{IsDel: constat.DELETE}).Error; err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
