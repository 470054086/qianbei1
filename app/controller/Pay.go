package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qianbei.com/app/params/request"
	"qianbei.com/app/params/response"
	"qianbei.com/app/service"
	"qianbei.com/util"
)

type Pay struct {
	payService service.PayService
}

func (p *Pay) Record(c *gin.Context) {
	r := &request.Record{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	if r.Order != "" && r.Desc != "" {
		if !(r.Order == "account" || r.Order == "created_at" ) {
			c.JSON(http.StatusInternalServerError, util.Error())
			return
		}
		if !(r.Desc == "desc" || r.Desc == "asc" ) {
			c.JSON(http.StatusInternalServerError, util.Error())
			return
		}
	}

	record, count, err := p.payService.Record(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	// 定义返回类型
	m := response.Record{
		Lists: record,
		Total: response.BaseTotal{
			Page:     r.Page,
			PageSize: r.PageSize,
			Count:    count,
		},
	}

	c.JSON(http.StatusOK, util.Success(m))
}

// 添加记录
func (p *Pay) AddRecord(c *gin.Context) {
	r := &request.PayRecord{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, util.ErrorParams())
		return
	}
	pay := p.payService.CreatePay(r.UploadType)
	record, err := pay.AddRecord(1, r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ErrorLog(err))
		return
	}
	c.JSON(http.StatusOK, util.Success(record))
}
