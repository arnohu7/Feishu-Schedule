package main

// formulate 函数用于格式化班表数据，便于后续的提交
func formulate() [4*int(NumSlots) - 2]map[string]interface{} {
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
