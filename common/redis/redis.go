package redis

import (
	"context"
	"time"

	"github.com/act-gpt/marino/config/system"

	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client
var ctx = context.Background()

const SessionKey = "SESSION:%s"

func hourSeconds(t time.Time) time.Duration {
	year, month, day := t.Date()
	hour := t.Hour()
	t2 := time.Date(year, month, day, hour, 0, 0, 0, t.Location())
	return t.Sub(t2)
}

func daySeconds(t time.Time) time.Duration {
	year, month, day := t.Date()
	t2 := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return t.Sub(t2)
}

func weekSeconds(t time.Time) time.Duration {
	currentTime := time.Now()
	// 计算本周剩余时间
	weekday := int(currentTime.Weekday())
	daysUntilSunday := 7 - weekday
	return time.Duration(daysUntilSunday*24)*time.Hour - time.Duration(currentTime.Hour())*time.Hour - time.Duration(currentTime.Minute())*time.Minute - time.Duration(currentTime.Second())*time.Second
}

func monthSeconds(t time.Time) time.Duration {
	year, month, _ := t.Date()
	t2 := time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
	return t.Sub(t2)
}

func yearSeconds(t time.Time) time.Duration {
	year, _, _ := t.Date()
	t2 := time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
	return t.Sub(t2)
}

func Init() {
	opts := ParseRedisOption()
	if opts != nil {
		Rdb = redis.NewClient(opts)
	}
}

func SetIncr(key string, value string) {
	Rdb.Set(ctx, key, value, 0)
}

func IncrBy(key string) {
	Rdb.Incr(ctx, key)
}

func IncrHourBy(key string) {
	Rdb.Incr(ctx, key)
	Rdb.Expire(ctx, key, hourSeconds(time.Now()))
}

func IncrDayBy(key string) {
	Rdb.Incr(ctx, key)
	Rdb.Expire(ctx, key, daySeconds(time.Now()))
}

func IncrWeekBy(key string) {
	Rdb.Incr(ctx, key)
	Rdb.Expire(ctx, key, weekSeconds(time.Now()))
}

func IncrMonthBy(key string) {
	Rdb.Incr(ctx, key)
	Rdb.Expire(ctx, key, monthSeconds(time.Now()))
}

func IncrYearBy(key string) {
	Rdb.Incr(ctx, key)
	Rdb.Expire(ctx, key, yearSeconds(time.Now()))
}

func GetFloat(key string) float64 {
	val, err := Rdb.Get(ctx, key).Float64()
	if err != nil {
		return 0.0
	}
	return val
}

func Increment(key string, field string, incr int64) {
	Rdb.HIncrBy(ctx, key, field, incr)
}

func GetLimited(key string, field string) (float64, error) {
	return Rdb.HGet(ctx, key, field).Float64()
}

func HSet(key string, field string, value string) {
	Rdb.HSet(ctx, key, field, value)
}

func HMGet(key string, fields ...string) ([]interface{}, error) {
	return Rdb.HMGet(ctx, key, fields...).Result()
}

func Set(key string, value string, t time.Duration) {
	Rdb.Set(ctx, key, value, t)
}

func Get(key string) (string, error) {
	return Rdb.Get(ctx, key).Result()
}

func ParseRedisOption() *redis.Options {
	url := system.Config.Redis.DataSource
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return opt
}

func Close() {
	if Rdb == nil {
		return
	}
	err := Rdb.Close()
	if err != nil {
		return
	}
}
