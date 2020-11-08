package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qianbei.com/app/cache"
	"qianbei.com/app/params/request"
	"qianbei.com/util"
)

type Category struct {
	cateCache cache.Category
}

// 进行查询
func (cate *Category) Index(c *gin.Context) {
	data, _ := cate.cateCache.Index(1)
	c.JSON(http.StatusOK, util.Success(data))
}

func (cate *Category) Create(c *gin.Context) {
	r := &request.CreateCategory{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	create, err := cate.cateCache.Create(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(create))
}

// 修改操作
func (cate *Category) Update(c *gin.Context) {
	r := &request.UpdateCategory{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	userCategory, err := cate.cateCache.Update(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(userCategory))
}

func (cate *Category) Del(c *gin.Context) {
	r := &request.DelCategory{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	err := cate.cateCache.Del(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(nil))
}
