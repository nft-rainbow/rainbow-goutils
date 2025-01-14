package timeutils

import "time"

// timestamp -> date
//
// func TimestampToSmileShopDate(timestamp int64) string {

// }

var ssLocation *time.Location
var mysqlLocation *time.Location

func init() {
	locationString := "Asia/Phnom_Penh"
	location_, err := time.LoadLocation(locationString)
	if err != nil {
		panic(err)
	}
	ssLocation = location_

	mysqlLocation = time.Local
}

// dt: 20230901 (yyyyMMdd)
// tm: 090000 (HHmmss)
func SsDtToTimestamp(dt string, tm string) (int64, error) {
	// 合并日期和时间字符串
	datetimeStr := dt + " " + tm

	// 使用指定的时区解析合并后的字符串
	t, err := time.ParseInLocation("20060102 150405", datetimeStr, ssLocation)
	if err != nil {
		return 0, err
	}

	// 返回 Unix 时间戳
	return t.Unix(), nil
}

func TimestampToSsDtTm(timestamp int64) (string, string) {

	// 将时间戳转换为指定时区的时间
	t := time.Unix(timestamp, 0).In(ssLocation)

	// 格式化为日期和时间字符串
	dt := t.Format("20060102")
	tm := t.Format("150405")

	return dt, tm
}

func TimestampToMysqlTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).In(mysqlLocation)
}

func MysqlTimeToTimestamp(t time.Time) int64 {
	return t.In(mysqlLocation).Unix()
}

func ParseCambodiaTime(dt string, tm string) (time.Time, error) {
	return time.ParseInLocation("20060102150405", dt+tm, ssLocation)
}

func GetTodayStart() time.Time {
	return time.Now().In(mysqlLocation).Truncate(24 * time.Hour)
}
