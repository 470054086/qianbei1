package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"qianbei.com/app/params/request"
	"qianbei.com/app/service"
	"qianbei.com/util"
)

type Book struct {
	bookService service.BookService
}

// 创建账本
func (b *Book) Create(c *gin.Context) {
	r := &request.AddBook{}
	if err := c.ShouldBind(r); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	m, err := b.bookService.Add(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(m))
}

// 查找账本
func (b Book) Index(c *gin.Context) {
	list, err := b.bookService.GetList(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(list))
}

// 修改账本
func (b *Book) Update(c *gin.Context) {
	r := &request.UpdateBook{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	m, err := b.bookService.Update(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(m))
}

// 删除账本
func (b *Book) Delete(c *gin.Context) {
	r := &request.DeleteBook{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	err := b.bookService.DeleteBook(r.UserId, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(nil))
}

// 加入用户
func (b *Book) JoinUser(c *gin.Context) {
	r := &request.JoinUserRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	userId := r.UserId
	m, err := b.bookService.UserJoin(userId, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(m))
}
