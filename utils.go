package main

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type Config struct {
	APPID        string `json:"APP_ID"`
	APPSecret    string `json:"APP_Secret"`
	APPToken     string `json:"APP_Token"`
	AccessToken  string `json:"Access_Token"`
	ReadTableID  string `json:"Read_Table_ID"`
	WriteTableID string `json:"Write_Table_ID"`
}

// 读取 config
func loadConfig(filename string) (Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil && err != io.EOF {
		return config, err
	}

	return config, nil
}

func getWeekday(field string) Weekday {
	switch {
	case strings.Contains(field, "周一"):
		return Monday
	case strings.Contains(field, "周二"):
		return Tuesday
	case strings.Contains(field, "周三"):
		return Wednesday
	case strings.Contains(field, "周四"):
		return Thursday
	case strings.Contains(field, "周五"):
		return Friday
	case strings.Contains(field, "周六"):
		return Saturday
	case strings.Contains(field, "周日"):
		return Sunday
	default:
		return -1
	}
}

func getTimeSlot(timeSlotStr string) TimeSlot {
	switch timeSlotStr {
	case "9:00-10:00":
		return Slot1
	case "10:00-12:00":
		return Slot2
	case "13:30-16:10":
		return Slot3
	case "16:10-18:00":
		return Slot4
	case "19:00-21:00":
		return Slot5
	default:
		return -1
	}
}

// getTimeSlotIndex 根据字段名称和时间段字符串确定时间段索引
func getTimeSlotIndex(field string, timeSlotStr string) int {
	weekday := getWeekday(field)
	timeSlot := getTimeSlot(timeSlotStr)

	if weekday >= 0 && timeSlot >= 0 {
		return int(weekday)*int(NumSlots) + int(timeSlot)
	}

	return -1 // 返回-1表示未找到匹配的时间段
}
