package timeutils_test

import (
	"testing"
	"time"

	"github.com/nft-rainbow/rainbow-goutils/utils/timeutils"
)

func TestSsDtToTimestamp(t *testing.T) {
	tests := []struct {
		dt       string
		tm       string
		expected int64
	}{
		{"20230901", "090000", 1693533600},
		{"20231231", "235959", 1704041999},
	}

	for _, test := range tests {
		timestamp, err := timeutils.SsDtToTimestamp(test.dt, test.tm)
		if err != nil {
			t.Errorf("DtToTimestamp(%s, %s) returned an error: %v", test.dt, test.tm, err)
		}
		if timestamp != test.expected {
			t.Errorf("DtToTimestamp(%s, %s) = %d; expected %d", test.dt, test.tm, timestamp, test.expected)
		}
	}
}

func TestDtToTimestampInvalidInput(t *testing.T) {
	_, err := timeutils.SsDtToTimestamp("20230901", "250000")
	if err == nil {
		t.Error("Expected error for invalid time, but got nil")
	}

	_, err = timeutils.SsDtToTimestamp("20230932", "120000")
	if err == nil {
		t.Error("Expected error for invalid date, but got nil")
	}
}

func TestTimestampToDtTm(t *testing.T) {
	tests := []struct {
		timestamp    int64
		locationName string
		expectedDt   string
		expectedTm   string
	}{
		{1693533600, "Asia/Phnom_Penh", "20230901", "090000"}, // 2023-09-01 09:00:00 in Asia/Phnom_Penh
		{1704041999, "Asia/Phnom_Penh", "20231231", "235959"}, // 2023-12-31 23:59:59 in Asia/Phnom_Penh
	}

	for _, test := range tests {
		dt, tm := timeutils.TimestampToSsDtTm(test.timestamp)

		if dt != test.expectedDt || tm != test.expectedTm {
			t.Errorf("TimestampToDtTm(%d, %s) = (%s, %s); expected (%s, %s)", test.timestamp, test.locationName, dt, tm, test.expectedDt, test.expectedTm)
		}
	}
}

func TestTimestampToMysqlTime(t *testing.T) {
	timestamp := int64(1693533600) // 2023-09-01 10:00:00 UTC+8
	expectedMysqlTime := "2023-09-01 10:00:00"

	mysqlTime := timeutils.TimestampToMysqlTime(timestamp)

	if mysqlTime.Format("2006-01-02 15:04:05") != expectedMysqlTime {
		t.Errorf("TimestampToMysqlTime(%d) = %s; 期望 %s", timestamp, mysqlTime.Format("2006-01-02 15:04:05"), expectedMysqlTime)
	}
}

func TestMysqlTimeToTimestamp(t *testing.T) {
	mysqlTime := "2023-09-01 10:00:00"
	expectedTimestamp := int64(1693533600)

	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", mysqlTime, time.Local)
	if err != nil {
		t.Fatalf("解析MySQL时间字符串失败: %v", err)
	}

	timestamp := timeutils.MysqlTimeToTimestamp(parsedTime)

	if timestamp != expectedTimestamp {
		t.Errorf("MysqlTimeToTimestamp(%s) = %d; 期望 %d", mysqlTime, timestamp, expectedTimestamp)
	}
}
