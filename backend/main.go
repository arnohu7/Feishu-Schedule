package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
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

type Weekday int

const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
	NumDays
)

type TimeSlot int

const (
	Slot1 TimeSlot = iota
	Slot2
	Slot3
	Slot4
	Slot5
	NumSlots
)

type Assignment struct {
	Name  string
	Count int
}

// scheduleMap 用来存放助理提交的班表
var scheduleMap map[string][int(NumSlots) * int(NumDays)]bool

// timeSlotAssignments 用来存储每个时间段的助理分配情况
var timeSlotAssignments [int(NumSlots) * int(NumDays)][]string

// assistantAssignments 用来存储每个助理已经被分配的时间段数量
var assistantAssignments map[string]int

func assign() {
	assistantAssignments = make(map[string]int)

	// 循环遍历每个时间段
	for i := 0; i < int(NumSlots)*int(NumDays); i++ {
		availableAssistants := []Assignment{}

		// 找出所有可以在这个时间段工作的助理
		for assistant, availability := range scheduleMap {
			if availability[i] {
				availableAssistants = append(availableAssistants, Assignment{
					Name:  assistant,
					Count: assistantAssignments[assistant],
				})
			}
		}

		// 按照助理的已被分配时间段的数量从小到大排序
		sort.Slice(availableAssistants, func(a, b int) bool {
			return availableAssistants[a].Count < availableAssistants[b].Count
		})

		// 选择最少被分配时间段的助理分配到这个时间段
		maxAssignments := 4
		if i%5 == 4 {
			// 对应 19:00-21:00 时间段
			maxAssignments = 2
		}
		for j := 0; j < len(availableAssistants) && j < maxAssignments; j++ {
			timeSlotAssignments[i] = append(timeSlotAssignments[i], availableAssistants[j].Name)
			assistantAssignments[availableAssistants[j].Name]++
		}
	}
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

// generateSchedule 函数用于生成最后的班表
func generateSchedule() [4*int(NumSlots) - 2]map[string]interface{} {
	var schedule [18]map[string]interface{}

	for i := range schedule {
		schedule[i] = make(map[string]interface{})
	}

	for i, assistants := range timeSlotAssignments {
		timeIndex := i % 5
		dayIndex := i / 5
		var day string
		var timeSlot string

		// 根据 dayIndex 来确定周几
		switch dayIndex {
		case 0:
			day = "周一"
		case 1:
			day = "周二"
		case 2:
			day = "周三"
		case 3:
			day = "周四"
		case 4:
			day = "周五"
		case 5:
			day = "周六"
		case 6:
			day = "周日"
		}

		// 根据timeIndex确定时间段
		switch timeIndex {
		case 0:
			timeSlot = "9:00-10:00"
		case 1:
			timeSlot = "10:00-12:00"
		case 2:
			timeSlot = "13:30-16:10"
		case 3:
			timeSlot = "16:10-18:00"
		case 4:
			timeSlot = "19:00-21:00"
		}

		// 计算写入位置的起始索引
		baseIndex := timeIndex * 4

		// 将助理数据写入 schedule
		for j, assistant := range assistants {
			scheduleIndex := baseIndex + j
			schedule[scheduleIndex][day] = assistant[1 : len(assistant)-1]
			schedule[scheduleIndex]["时间"] = timeSlot
		}
	}
	return schedule
}

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 创建 Client
	client := lark.NewClient(config.APPID, config.APPSecret)
	// 创建请求对象
	req := larkbitable.NewListAppTableRecordReqBuilder().
		AppToken(config.APPToken).
		TableId(config.ReadTableID).
		Build()

	// 发起请求
	resp, err := client.Bitable.AppTableRecord.List(context.Background(), req, larkcore.WithUserAccessToken(config.AccessToken))

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return
	}

	// 业务处理
	// 创建一个字典来存储每个人的姓名以及可值班的时间
	scheduleMap = make(map[string][int(NumSlots) * int(NumDays)]bool)

	for _, item := range resp.Data.Items {
		// 从"提交人"map中获取名字
		submitterMap, ok := item.Fields["提交人"].(map[string]interface{})
		if !ok {
			continue // 如果第一个元素不是 map，则跳过此纪录
		}

		submitter, ok := submitterMap["name"].(string)
		if !ok {
			continue // 如果名字字段不存在或类型不正确，则跳过此记录
		}

		// 初始化一个布尔数组来存储可值班时间
		var schedule [int(NumSlots) * int(NumDays)]bool

		// 遍历每个字段，查找与可值班时间相关的字段
		for field, value := range item.Fields {
			if strings.Contains(field, "可值班时间") {
				timeSlots, ok := value.([]interface{})
				if ok {
					for _, timeSlot := range timeSlots {
						timeSlotStr, ok := timeSlot.(string)
						if ok {
							timeSlotIndex := getTimeSlotIndex(field, timeSlotStr)
							if timeSlotIndex >= 0 {
								schedule[timeSlotIndex] = true
							}
						}
					}
				}
			}
		}

		// 将提交人和可值班时间添加到字典中
		scheduleMap[larkcore.Prettify(submitter)] = schedule
	}

	// 调用 assign 函数进行排班
	assign()

	// 查看排班结果
	for i, assistants := range timeSlotAssignments {
		fmt.Printf("Time Slot %d: %v\n", i, assistants)
	}

	// 生成格式化的排班结果
	schedule := generateSchedule()

	// 查看格式化的排班结果
	for i, item := range schedule {
		fmt.Printf("schedule[%d]: %v\n", i, item)
	}

	// 遍历 schedule 数组,为每条记录创建并发送请求
	for _, record := range schedule {
		req_put := larkbitable.NewCreateAppTableRecordReqBuilder().
			AppToken(config.APPToken).
			TableId(config.WriteTableID).
			AppTableRecord(larkbitable.NewAppTableRecordBuilder().
				Fields(record).
				Build()).
			Build()

		resp, err := client.Bitable.AppTableRecord.Create(context.Background(), req_put, larkcore.WithUserAccessToken(config.AccessToken))

		if err != nil {
			// 打印错误并可能终止循环
			fmt.Println("Error:", err)
			break
		}

		if !resp.Success() {
			// 处理服务器返回的错误
			fmt.Println("Server Error:", resp.Code, resp.Msg, resp.RequestId())
		}
	}
}
