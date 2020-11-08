package response

import "qianbei.com/app/model"

type Category struct {
	Pay    []*model.Category `json:"pay"`
	Income []*model.Category `json:"income"`
}
