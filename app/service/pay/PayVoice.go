package pay

import (
	"qianbei.com/app/params/request"
	"qianbei.com/app/params/response"
)

type payVoice struct {
}

func (p payVoice) AddRecord(userId int, r *request.PayRecord) (*response.AddRecord, error) {
	panic("implement me")
}

func (p payVoice) UpdateRecord() {
	panic("implement me")
}

func (p payVoice) DeleteRecord() {
	panic("implement me")
}

func NewPayVoice() *payVoice {
	return &payVoice{}
}
