package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkbitable "github.com/larksuite/oapi-sdk-go/v3/service/bitable/v1"
)

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
		submitter = larkcore.Prettify(submitter)
		scheduleMap[submitter[1:len(submitter)-1]] = schedule
	}

	// 调用 assign 函数进行排班
	assign()

	// 查看排班结果
	for i, assistants := range timeSlotAssignments {
		fmt.Printf("Time Slot %d: %v\n", i, assistants)
	}

	// 生成格式化的排班结果
	schedule := formulate()

	// 查看格式化的排班结果
	for i, item := range schedule {
		fmt.Printf("schedule[%d]: %v\n", i, item)
	}

	// 遍历 schedule 数组,为每条记录创建并发送请求
	for i, record := range schedule {
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
			fmt.Printf("第 %d 条记录提交出错", i)
			fmt.Println("Server Error:", resp.Code, resp.Msg, resp.RequestId())
		}

		if i%4 == 3 {
			req_put_blank := larkbitable.NewCreateAppTableRecordReqBuilder().
				AppToken(config.APPToken).
				TableId(config.WriteTableID).
				AppTableRecord(larkbitable.NewAppTableRecordBuilder().
					Fields(map[string]interface{}{`时间`: ``, `周一`: ``, `周二`: ``, `周三`: ``, `周四`: ``, `周五`: ``, `周六`: ``, `周日`: ``}).
					Build()).
				Build()

			resp, err := client.Bitable.AppTableRecord.Create(context.Background(), req_put_blank, larkcore.WithUserAccessToken(config.AccessToken))

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
}
