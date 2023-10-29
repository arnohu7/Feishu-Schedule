package main

import "sort"

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
