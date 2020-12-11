package utils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//时间类型
type TimeType int

const (
	YearType   TimeType = 1
	MonthType  TimeType = 2
	DayType    TimeType = 3
	HourType   TimeType = 4
	MinuteType TimeType = 5
	SecondType TimeType = 6
)

//时间字符串转换成时间
func ParseStrToTime(timeStr string, flag int) time.Time {
	var t time.Time
	var err error
	if flag == 1 {
		t, err = time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	} else if flag == 2 {
		t, err = time.ParseInLocation("2006-01-02 15:04", timeStr, time.Local)
	} else if flag == 3 {
		t, err = time.ParseInLocation("2006-01-02", timeStr, time.Local)
	} else if flag == 4 {
		t, err = time.ParseInLocation("2006.01.02 15:04:05", timeStr, time.Local)
	} else if flag == 5 {
		t, err = time.ParseInLocation("2006.01.02 15:04", timeStr, time.Local)
	} else {
		t, err = time.ParseInLocation("2006.01.02", timeStr, time.Local)
	}
	if err != nil {
		log.Println("convert [" + timeStr + "] string to time is failed")
		t, _ = time.ParseInLocation("2006-01-02", "1000-01-01", time.Local)
	}
	return t
}

/**
  获取多少天,多少月或者多少年之前或之后的时间
  dayRange: 间隔的天数，月数或者年份数
  timeType: 决定是取天数，月数还是年数
*/
func GetAfterDay(dayRange int, timeType TimeType) time.Time {
	now := time.Now()
	var tmpDay time.Time
	if timeType == YearType {
		tmpDay = now.AddDate(dayRange, 0, 0)
	} else if timeType == MonthType {
		tmpDay = now.AddDate(0, dayRange, 0)
	} else if timeType == DayType {
		tmpDay = now.AddDate(0, 0, dayRange)
	} else if timeType == HourType {
		tmpDay = GetAfterTime(strconv.Itoa(dayRange) + "h")
	} else if timeType == MinuteType {
		tmpDay = GetAfterTime(strconv.Itoa(dayRange) + "m")
	} else if timeType == SecondType {
		tmpDay = GetAfterTime(strconv.Itoa(dayRange) + "s")
	} else {
		tmpDay = now.AddDate(0, 0, dayRange)
	}
	return tmpDay
}

/**
  获取多少小时，分钟及秒之前或之后的时间
  timeRange: 时间差，比如：
      10h     获取10小时之后的时间
      -10h    获取10小时之前的时间
      10m     获取10分钟之后的时间
      -10m    获取10分钟之后的时间
      10s     获取10秒之后的时间
      -10s    获取10秒之后的时间
*/
func GetAfterTime(timeRange string) time.Time {
	m, _ := time.ParseDuration(timeRange)
	tmp := time.Now().Add(m)
	return tmp
}

/**
  把 datetime 转换成 时间字符串
  t: datetime 时间，比如：2019-09-17 09:45:42.5962359 +0800 CST m=+0.003989201
  flag: 标识位，决定输出的时间字符串的格式
*/
func ParseTimeToStr(t time.Time, flag int) string {
	var timeStr string
	if flag == 1 {
		timeStr = t.Format("2006-01-02 15:04:05")
	} else if flag == 2 {
		timeStr = t.Format("2006-01-02 15:04")
	} else if flag == 3 {
		timeStr = t.Format("2006-01-02")
	} else if flag == 4 {
		timeStr = t.Format("2006.01.02 15:04:05")
	} else if flag == 6 {
		timeStr = t.Format("2006.01.02 15:04")
	} else {
		timeStr = t.Format("2006.01.02")
	}
	return timeStr
}

/**
  把 datetime 转换成时间戳
  t: datetime 时间
*/
func ParseTimeToInt64(t time.Time) int64 {
	return t.Unix()
}

//时间戳转换成年月日
func ParseTimeToDate(timeStr string) string {
	dateTime := time.Unix(time.Now().Unix(), 0).Format(timeStr)
	return dateTime
}

//时间转换 将1993-12-26 10:30:00转换为time
func ParseTimeByTimeStr(str, errPrefix string) (time.Time, error) {
	p := strings.TrimSpace(str)
	if p == "" {
		return time.Time{}, errors.New(fmt.Sprintf("%s不能为空", errPrefix))
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return time.Time{}, errors.New(fmt.Sprintf("%s格式错误", errPrefix))
	}

	return t, nil
}

//获取最近的周一
func ParseCurrentMonday(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStart
}

//返回某一天的当地时区0点
func ParseMorningTime(t time.Time) time.Time {
	s := t.Format("19931226")
	result, _ := time.ParseInLocation("19931226", s, time.Local)
	return result
}

//当月第一天0点
func ParseFirstDayOfMonthMorning(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

//获取传入时间前一天的时间，不传默认是昨天
func ParseYesterdayTime(t ...time.Time) time.Time {
	if len(t) == 0 {
		return time.Now().AddDate(0, 0, -1)
	} else {
		return t[0].AddDate(0, 0, -1)
	}
}

//把int64转换成1993-12-26 10:30:00
func ParseTimeToTimeStr(intTime int64, strfmt ...string) string {
	t := time.Unix(intTime/1e3, 0)
	defaultFmt := "2006-01-02 15:04:05"
	if len(strfmt) > 0 {
		defaultFmt = strfmt[0]
	}
	return t.Format(defaultFmt)
}

func DateFormat(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}

//判断是否为同一天
func IsSameDate(date1 time.Time, date2 time.Time) bool {
	if date1.YearDay() == date2.YearDay() {
		return true
	}
	return false
}

// 计算日期相差多少天
func SubDays(t1, t2 time.Time) (day int) {
	day = int(t1.Sub(t2).Hours() / 24)
	return
}

// 计算日期相差多少小时
func SubHours(t1, t2 time.Time) (hour int) {
	hour = int(t1.Sub(t2).Hours())
	return
}

// 返回星期
func GetWeekName(week int32) string {
	if week >= 0 && week < 7 {
		return []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}[week]
	}

	return ""
}