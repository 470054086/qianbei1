package util

import (
	"github.com/Gre-Z/common/jtime"
	"qianbei.com/constat"
	"strconv"
	"time"
)

//获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

//获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回下个月的零点
func GetStartDataOfNextMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, 0)
}

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取年月日的显示
func GetYmShow(d jtime.JsonTime) string {
	year := time.Unix(d.Unix(), 0).Year()
	month := constat.Month[time.Unix(d.Unix(), 0).Month().String()]
	// 将year转化为字符串
	return strconv.Itoa(year) + strconv.Itoa(month)
}

// 获取本月的年月显示
func GetCurrYmShow() string {
	year := time.Now().Year()
	month := constat.Month[time.Now().Month().String()]
	return strconv.Itoa(year) + strconv.Itoa(month)
}
