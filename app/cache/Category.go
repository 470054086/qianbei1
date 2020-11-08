package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"qianbei.com/app/model"
	"qianbei.com/app/params/request"
	"qianbei.com/constat"
	"qianbei.com/core"
	"sort"
	"time"
)

const (
	UserCategoryMin = 10000 // 定义用户定义的分类最小ID 为了区分
)

type Category struct {
	cateModel     *model.Category
	cateUserModel *model.CategoryUser
}

func (c *Category) Index(userId int) (model.CategorySlice, error) {
	ctx := context.Background()
	redisKey := fmt.Sprintf(constat.REDIS_CATEGORY_USERID, userId)
	val := core.Redis().Get(ctx, redisKey).Val()
	if val != "" {
		var res model.CategorySlice
		_ = json.Unmarshal([]byte(val), &res)
		return res, nil
	}
	return c.createRedis(userId)
}

func (c *Category) Create(userId int, r *request.CreateCategory) (*model.CategoryUser, error) {
	err := c.isExistName(userId, r.Name, r.Types, *r.Pid)
	if err != nil {
		return nil, err
	}
	m := &model.CategoryUser{
		Pid:    *r.Pid,
		UserId: userId,
		Name:   r.Name,
		Types:  r.Types,
		Remark: "",
		Sorts:  *r.Sorts,
	}
	if err := core.Db().Save(&m).Error; err != nil {
		return nil, errors.Wrap(err, "")
	}
	go func(userId int) {
		c.createRedis(userId)
	}(userId)

	return m, nil
}

// 构造redis结构
func (c *Category) createRedis(userId int) (model.CategorySlice, error) {
	redisKey := fmt.Sprintf(constat.REDIS_CATEGORY_USERID, userId)
	// 获取全部的分类
	list, err := c.cateModel.GetList()
	// 获取用户的分类
	userList, err := c.cateUserModel.GetList(1)
	concat := c.concat(list, userList)
	// 排序
	categorySort := c.cateSorts(concat)
	// 数据转换 进行递归处理
	category := c.deep(categorySort, 0)
	// 进行redis的处理
	marshal, _ := json.Marshal(category)
	core.Redis().Set(context.Background(), redisKey, marshal, time.Hour*7*24)
	return category, err
}

func (c *Category) Update(userId int, r *request.UpdateCategory) (*model.CategoryUser, error) {
	// 如果小于的话 说明修改的原始分类
	var userCategory *model.CategoryUser
	var err error
	if r.Id < UserCategoryMin {
		// 如果不存在的话 则添加一条记录
		userCategory, err = c.cateUserModel.GetFirstByCid(r.Id, userId)
		if err != nil {
			return nil, err
		}
		// 说明已经存在 直接修改即可
		if userCategory != nil {
			// 如果已经被删除了
			if userCategory.IsDel == constat.DELETE {
				return nil,constat.NewShowMessage(constat.CATEGORY_NOT_EXIST)
			}
			// 如果修改的是name 需要判断name是否存在
			if r.Name != "" {
				err := c.isExistName(userId, r.Name, userCategory.Types, userCategory.Pid)
				if err != nil {
					return nil, err
				}
				userCategory.Name = r.Name
			}
			if r.Sorts != nil {
				userCategory.Sorts = *r.Sorts
			}
			if err := core.Db().Save(userCategory).Error; err != nil {
				return nil, errors.Wrap(err, "")
			}
		} else {
			// 获取原来的pid和id
			cate, err := c.cateModel.GetFirstById(r.Id)
			if err != nil {
				return nil, err
			}
			if cate == nil {
				return nil, constat.NewShowMessage(constat.CATEGORY_NOT_EXIST)
			}
			// 如果已经存在的话 直接修改
			if userCategory == nil {
				// 添加操作
				m := &model.CategoryUser{
					UserId: userId,
					Pid:    cate.Pid,
					Cid:    cate.ID,
					Types:  cate.Types,
					Remark: cate.Remark,
				}
				if r.Name != "" {
					m.Name = r.Name
				} else {
					m.Name = cate.Name
				}
				if r.Sorts != nil {
					m.Sorts = *r.Sorts
				} else {
					m.Sorts = cate.Sorts
				}

				if err := core.Db().Save(m).Error; err != nil {
					return nil, errors.Wrap(err, "")
				}
			} else {
				if r.Name == "" {
					userCategory.Name = r.Name
				}
				if r.Sorts != nil {
					userCategory.Sorts = *r.Sorts
				}
				if err := core.Db().Save(userCategory).Error; err != nil {
					return nil, errors.Wrap(err, "")
				}
			}

		}
	} else {
		// 如果是userCategory 直接修改
		userCategory, err = c.cateUserModel.GetFirstById(r.Id)
		if err != nil {
			return nil, err
		}
		if userCategory == nil {
			return nil, constat.NewShowMessage(constat.CATEGORY_NOT_EXIST)
		}
		if r.Name != "" {
			err := c.isExistName(userId, r.Name, userCategory.Types, userCategory.Pid)
			if err != nil {
				return nil, err
			}
			userCategory.Name = r.Name
		}
		if r.Sorts != nil {
			userCategory.Sorts = *r.Sorts
		}
		if err := core.Db().Save(userCategory).Error; err != nil {
			return nil, errors.Wrap(err, "")
		}
	}
	// 构建redis缓存
	go func(userId int) {
		c.createRedis(userId)
	}(userId)

	return userCategory, nil
}

func (c *Category) Del(userId int, r *request.DelCategory) error {
	// 小于的话 说明是删除的系统类型
	if r.Id < UserCategoryMin {
		// 判断系统类型是否存在
		category, err := c.isExistByCidCategory(r.Id)
		if err != nil {
			return err
		}
		if category == nil {
			// 添加数据
			m := &model.CategoryUser{
				UserId: userId,
				Pid:    category.Pid,
				Cid:    category.ID,
				Name:   category.Name,
				Types:  category.Types,
				Remark: category.Remark,
				Sorts:  category.Sorts,
				IsDel:  constat.DELETE,
			}
			if err := core.Db().Save(m).Error; err != nil {
				return errors.Wrap(err, "")
			}
		}
		category.IsDel = constat.DELETE
		if err := core.Db().Save(category).Error; err != nil {
			return errors.Wrap(err, "")
		}
	} else {
		user, err := c.isExistCategoryUser(r.Id)
		if err != nil {
			return errors.Wrap(err, "")
		}
		user.IsDel = constat.DELETE
		if err := core.Db().Save(user).Error; err != nil {
			return errors.Wrap(err, "")
		}
	}
	// 构建redis缓存
	go func(userId int) {
		c.createRedis(userId)
	}(userId)
	return nil
}

func (c *Category) isExistCategory(id int) (*model.Category, error) {
	cate, err := c.cateModel.GetFirstById(id)
	if err != nil {
		return nil, err
	}
	if cate == nil {
		return nil, constat.NewShowMessage(constat.CATEGORY_NOT_EXIST)
	}
	return cate, nil
}
func (c *Category) isExistByCidCategory(cid int) (*model.CategoryUser, error) {
	cate, err := c.cateUserModel.GetFirstByCId(cid)
	if err != nil {
		return nil, err
	}
	if cate == nil {
		return nil, constat.NewShowMessage(constat.CATEGORY_NOT_EXIST)
	}
	return cate, nil
}

func (c *Category) isExistCategoryUser(id int) (*model.CategoryUser, error) {
	cate, err := c.cateUserModel.GetFirstById(id)
	if err != nil {
		return nil, err
	}
	if cate == nil {
		return nil, constat.NewShowMessage(constat.CATEGORY_NOT_EXIST)
	}
	return cate, nil
}

// 判断名称是否存在
func (c *Category) isExistName(userId int, name string, types int, pid int) error {
	// 判断系统类型的分类下是否存在这个名字
	cate, err := c.cateModel.GetFirstByNameAndPid(name, pid, types)
	if err != nil {
		return err
	}
	if cate != nil {
		// 判断自定义下面是否存在删除这个名字
		cateUser, err := c.cateUserModel.GetFirstByCNameAndUserIdAndPidDel(name, userId, pid, types)
		if err != nil {
			return err
		}
		if cateUser == nil {
			return constat.NewShowMessage(constat.CATEGORY_NAME_EXIST)
		}
	}
	// 判断自定义添加的是否存在这个名字
	cateUser, err := c.cateUserModel.GetFirstByCNameAndUserIdAndPid(name, userId, pid, types)
	if err != nil {
		return err
	}
	if cateUser != nil {
		return constat.NewShowMessage(constat.CATEGORY_NAME_EXIST)
	}
	return nil
}

// 组合数据
func (c *Category) concat(list []*model.Category, userList []*model.CategoryUser) []*model.Category {
	if userList == nil {
		return list
	}
	res := []*model.Category{}
	// 进行循环分组
	for _, v := range list {
		for _, v1 := range userList {
			// 说明是删除系统自带的分类
			if v1.Cid != 0 && v.ID == v1.Cid {
				v.Isdel = v1.IsDel
				v.Sorts = v1.Sorts
				v.Name = v1.Name
			}
		}
		//  组合category的原生数据
		res = append(res, v)
	}

	// categoryUser进行组合
	for _, v := range userList {
		if v.Cid == 0 {
			t := &model.Category{
				ID:        v.ID,
				Pid:       v.Pid,
				Name:      v.Name,
				Types:     v.Types,
				Remark:    v.Remark,
				Sorts:     v.Sorts,
				Isdel:     v.IsDel,
				CreatedAt: v.CreatedAt,
				List:      nil,
			}
			res = append(res, t)
		}
	}
	return res
}

// 进行排序
func (c *Category) cateSorts(res []*model.Category) model.CategorySlice {
	var categorySort model.CategorySlice
	categorySort = append(categorySort, res...)
	sort.Sort(categorySort)
	return categorySort
}

func (c *Category) deep(list []*model.Category, pid int) []*model.Category {
	treeList := []*model.Category{}
	for _, v := range list {
		if v.Pid == pid {
			child := c.deep(list, v.ID)
			node := v
			node.List = child
			treeList = append(treeList, node)
		}
	}
	return treeList
}
