package main

import "sort"

type Assignment struct {
	Name  string
	Count int
}

// 预处理
func preProcess() {
	assistantAssignments = make(map[string]int)
}

// 将小黑屋的负责人进行排班，SDR 是 Small-Dark-Room 的缩写
// TODO: 使用非硬编码的方式实现这部分的逻辑
func assignSDRPricipal() {
	timeSlotAssignments[4] = append(timeSlotAssignments[4], "黄海洋")
	timeSlotAssignments[9] = append(timeSlotAssignments[9], "胡泽钊")
	timeSlotAssignments[14] = append(timeSlotAssignments[14], "宋彦斌")
	timeSlotAssignments[19] = append(timeSlotAssignments[19], "苏梓玲")
	timeSlotAssignments[24] = append(timeSlotAssignments[24], "杨毅")
	timeSlotAssignments[29] = append(timeSlotAssignments[29], "胡泽钊")
	timeSlotAssignments[34] = append(timeSlotAssignments[34], "杨毅")
}

// 对前台负责人进行排班，R 是 Reception 的缩写
// TODO: 使用非硬编码的方式实现这部分的逻辑
func assignRPrincipal() {

}

// 分配普通助理
func assignNormalAssistant() {
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

func assign() {
	preProcess()
	assignSDRPricipal()
	assignNormalAssistant()
}
