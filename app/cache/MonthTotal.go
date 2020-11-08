package cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"qianbei.com/constat"
	"qianbei.com/core"
	"qianbei.com/util"
	"strconv"
	"time"
)

type MonthTotalParams struct {
	UserId int `json:"-"`
	BookId int `json:"-"`
	Pay    int `json:"pay"`
	Income int `json:"income"`
}

type MonthTotal struct {
}

// 增加本月支出 缓存处理 月末过期 从数据库入库
func (m MonthTotal) AddTotal(userId int, account int, bookId int, types int) (*MonthTotalParams, error) {
	redisKey := fmt.Sprintf(constat.REDIS_MONTH_TOTAL, userId, bookId)
	ctx := context.Background()
	// 如果是支出的话
	if types == constat.PAY_EXPEND {
		// 如果不存在的话 直接添加并且设立过期时间
		if exist := core.Redis().Exists(ctx, redisKey).Val(); exist == 0 {
			m := map[string]interface{}{
				"user_id": userId,
				"pay":     account,
				"book_id": bookId,
				"income":  0,
			}
			if err := core.Redis().HMSet(ctx, redisKey, m).Err(); err != nil {
				return nil, err
			}
			expire := time.Second * time.Duration(util.GetStartDataOfNextMonth(time.Now()).Unix()-time.Now().Unix())
			if err := core.Redis().Expire(ctx, redisKey, expire).Err(); err != nil {
				return nil, err
			}
		} else {
			if err := core.Redis().HIncrBy(ctx, redisKey, "pay", int64(account)).Err(); err != nil {
				return nil, err
			}
		}

	} else if types == constat.PAY_INCOME {
		// 如果不存在的话 直接添加并且设立过期时间
		if exist := core.Redis().HGetAll(ctx, redisKey); exist == nil {
			m := map[string]interface{}{
				"user_id": userId,
				"pay":     0,
				"income":  account,
			}
			if err := core.Redis().HMSet(ctx, redisKey, m).Err(); err != nil {
				return nil, err
			}
			expire := time.Duration(util.GetStartDataOfNextMonth(time.Now()).Unix())
			if err := core.Redis().Expire(ctx, redisKey, expire).Err(); err != nil {
				return nil, err
			}
		} else {
			if err := core.Redis().HIncrBy(ctx, redisKey, "income", int64(account)).Err(); err != nil {
				return nil, err
			}
		}
	}
	result, _ := core.Redis().HGetAll(ctx, redisKey).Result()
	res := &MonthTotalParams{}
	// todo 稍后在处理这个问题
	for k, v := range result {
		if k == "pay" {
			atoi, _ := strconv.Atoi(v)
			res.Pay = atoi
		} else if k == "income" {
			atoi, _ := strconv.Atoi(v)
			res.Income = atoi
		} else if k == "user_id" {
			atoi, _ := strconv.Atoi(v)
			res.UserId = atoi
		} else if k == "book_id" {
			atoi, _ := strconv.Atoi(v)
			res.BookId = atoi
		}
	}
	return res, nil
}

func (m MonthTotal) DecrPay(userId int, bookId int, account int) error {
	redisKey := fmt.Sprintf(constat.REDIS_MONTH_TOTAL, userId, bookId)
	if err := core.Redis().HIncrBy(context.Background(), redisKey, "pay", -int64(account)).Err(); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func (m MonthTotal) DecrIncome(userId int, bookId int, account int) error {
	redisKey := fmt.Sprintf(constat.REDIS_MONTH_TOTAL, userId, bookId)
	if err := core.Redis().HIncrBy(context.Background(), redisKey, "income", -int64(account)).Err(); err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
