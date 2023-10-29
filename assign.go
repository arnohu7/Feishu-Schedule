package main

import (
	"fmt"
	"sort"
)

type Assignment struct {
	Name       string
	TotalHours int
}

type RPrincipleSet map[string]bool

var rPrincipleSet RPrincipleSet

// Add 方法用于将一个切片的元素添加到集合中
func (set RPrincipleSet) Add(items ...string) {
	for _, item := range items {
		set[item] = true
	}
}

func slotDuration(i int) int {
	switch i % 5 {
	case 0:
		return 60
	case 1, 4:
		return 120
	case 2:
		return 160
	case 3:
		return 110
	}
	return 0
}

// 预处理
// TODO: 使用非硬编码的方式实现这部分的逻辑
func preProcess() {
	assistantAssignments = make(map[string]int)

	RPrincipleToAdd := []string{
		"黄宝莹",
		"黄洪彬",
		"黄培轩",
		"赖广麟",
		"林亮秋",
		"罗雪源",
		"裴江博",
		"秦绍润",
		"宋彦斌",
		"苏梓玲",
		"唐苛耕",
		"许泽杭",
		"杨毅",
		"叶桂昂",
		"赵鑫",
		"郑桐",
	}
	rPrincipleSet = make(RPrincipleSet)
	rPrincipleSet.Add(RPrincipleToAdd...)
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

// 分配助理
func assignAssistant() {
	// 循环遍历每个时间段
	for i := 0; i < int(NumSlots)*int(NumDays); i++ {
		availableAssistants := []Assignment{}
		rPrincipalAssigned := len(timeSlotAssignments[i]) != 0

		// 找出当前时间段所有可以工作的助理
		for assistant, availability := range scheduleMap {
			if availability[i] {
				totalHours := assistantAssignments[assistant]
				availableAssistants = append(availableAssistants, Assignment{
					Name:       assistant,
					TotalHours: totalHours,
				})
			}
		}

		// 如果当前时间段是 9:00-10:00 ，让可以在两个时间段同时值班的助理优先
		// 其它时间段则只考虑目前已经安排的工作时长
		if i%int(NumSlots) == 0 {
			sort.SliceStable(availableAssistants, func(a, b int) bool {
				aAvailableBoth := scheduleMap[availableAssistants[a].Name][i] && scheduleMap[availableAssistants[a].Name][i+1]
				bAvailableBoth := scheduleMap[availableAssistants[b].Name][i] && scheduleMap[availableAssistants[b].Name][i+1]

				if aAvailableBoth == bAvailableBoth {
					return availableAssistants[a].TotalHours < availableAssistants[b].TotalHours
				} else {
					return aAvailableBoth && !bAvailableBoth
				}
			})

		} else {
			sort.Slice(availableAssistants, func(a, b int) bool {
				return availableAssistants[a].TotalHours < availableAssistants[b].TotalHours
			})
		}

		if i%5 != 4 && !rPrincipalAssigned {
			for j := 0; j < len(availableAssistants); j++ {
				if _, isRPrincipal := rPrincipleSet[availableAssistants[j].Name]; isRPrincipal {
					timeSlotAssignments[i] = append(timeSlotAssignments[i], availableAssistants[j].Name)
					assistantAssignments[availableAssistants[j].Name] += slotDuration(i)

					// 9:00-10:00 特殊判断，若某个负责人已经在 9:00-10:00 上班了且下个时间有空，则他应该也要在 10:00-12:00 上班
					if i%5 == 0 {
						if scheduleMap[availableAssistants[j].Name][i+1] {
							timeSlotAssignments[i+1] = append(timeSlotAssignments[i+1], availableAssistants[j].Name)
							assistantAssignments[availableAssistants[j].Name] += slotDuration(i + 1)
						}
					}

					rPrincipalAssigned = true
					break
				}
			}
		}

		if !rPrincipalAssigned && i%5 != 4 {
			fmt.Println("警告：没有可用的负责人在当前的时段", i)
		}

		// 对其它助理进行排班
		maxAssignments := 4
		if i%5 == 4 {
			maxAssignments = 3
		}
		for j := 0; j < len(availableAssistants) && len(timeSlotAssignments[i]) < maxAssignments; j++ {
			alreadyAssigned := false
			for _, assignedAssistant := range timeSlotAssignments[i] {
				if assignedAssistant == availableAssistants[j].Name {
					alreadyAssigned = true
					break
				}
			}
			if !alreadyAssigned {
				timeSlotAssignments[i] = append(timeSlotAssignments[i], availableAssistants[j].Name)
				assistantAssignments[availableAssistants[j].Name] += slotDuration(i)

				// 9:00-10:00 特殊判断，若某个助理已经在 9:00-10:00 上班了且下个时间有空，则他应该也要在 10:00-12:00 上班
				if i%5 == 0 {
					if scheduleMap[availableAssistants[j].Name][i+1] {
						timeSlotAssignments[i+1] = append(timeSlotAssignments[i+1], availableAssistants[j].Name)
						assistantAssignments[availableAssistants[j].Name] += slotDuration(i + 1)
					}
				}
			}
		}
	}
}

func assign() {
	preProcess()
	assignSDRPricipal()
	assignAssistant()
}
