package time

import (
	"github.com/gogo/protobuf/sortkeys"
	"github.com/pkg/errors"
	"time"
)

const DATE_BASE_FORMAT = "2006-01-02"							//日期格式字符串

const TIME_BASE_FORMAT string = "2006-01-02 15:04:05"					//时间格式字符串
const TIME_TAND_FORMAT string = "2006-01-02T15:04:05"
const TIME_MSEC_FORMAT string = "2006-01-02 15:04:05.000"

const TIME_DB_TIMESTAMP_FORMAT = "2006-01-02T15:04:05+08:00"

/**
 * 字符串转换为时间
 */
func ParseTime(str string, format string) (t time.Time,err error) {
	switch format {
	case "TB" :
		return time.Parse(TIME_BASE_FORMAT, str)
	case "TM" :
		return time.Parse(TIME_MSEC_FORMAT, str)
	case "DT" :
		return time.Parse(TIME_DB_TIMESTAMP_FORMAT, str)
	default :
		return time.Now(),errors.New("parse not support")
	}
}

/**
 * 当前时间毫秒数
 */
 func UnixMilli() int64 {
	 return time.Now().UnixNano() / 1000000
 }

/**
 * 时间切片函数
 */
func UnixTimeSlice(lTimeUnix []int64, iTimeSecond int64) []map[string]int64 {
	if iTimeSecond < 60 {
		//iTimeSecond 必须大于 60
		iTimeSecond = 60
	}

	var mUnixTime []map[string]int64
	if len(lTimeUnix) <= 0 {
		// 若输入时间为空，则返回同样为空
		return mUnixTime
	}

	// 时间戳排序
	sortkeys.Int64s(lTimeUnix)

	// 开始时间向前30秒
	begin_time := lTimeUnix[0] - 30
	end_time := begin_time + iTimeSecond
	mUnixTime = append(mUnixTime,map[string]int64{"begin_time" : begin_time, "end_time" : end_time})

	for _,value := range lTimeUnix {
		if value >= end_time {
			begin_time = value - 30
			end_time = begin_time + iTimeSecond
			mUnixTime = append(mUnixTime,map[string]int64{"begin_time" : begin_time, "end_time" : end_time})
		}
	}

	return mUnixTime
}