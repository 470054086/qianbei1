package request

import (
	"mime/multipart"
)

// 记账APP请求接口
type PayRecord struct {
	Category   int                   `json:"category_id"`                    // 分类
	Account    int                   `json:"account" binding:"required"`     // 金额
	BookId     int                   `json:"book_id" binding:"required" `    // 账本id
	Platform   int                   `json:"platform"  binding:"required"`   // 上传平台
	Types      int                   `json:"types" binding:"required"`       // 收入类型 1 支出 2 收入
	AssetsId   int                   `json:"assets_id"`                      // 资产id
	UploadType int                   `json:"upload_type" binding:"required"` // 上传类型 1 正常 2 语音 3 多条上传 4 图片 5 excel
	Memo       string                `json:"memo"`                           // 备注 语音只需要上传备注
	Excel      *multipart.FileHeader `form:"upload_key"`                     // 上传文件
	Image      *multipart.FileHeader `form:"image_key"`                      // 上传图片
}

// 查询的提供
type Record struct {
	StartAccount int    `json:"start_account"`              // 最小金额
	EndAccount   int    `json:"end_account"`                //最大金额
	BookId       int    `json:"book_id" binding:"required"` //账本id
	Types        int    `json:"types"`                      // 类型
	StartTime    string `json:"start_time"`                 //开始时间
	EndTime      string `json:"end_time"`                   // 结束时间
	Memo         string `json:"memo"`                       //备注信息
	Order        string `json:"order"`                      // 排序字段 支持时间和金额
	Desc         string `json:"desc"`                       //  倒序还是正序
	BaseRequest
}
