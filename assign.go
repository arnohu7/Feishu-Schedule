package main

import (
	"fmt"
	"sort"
)

type Assignment struct {
	Name  string
	Count int
}

type RPrincipleSet map[string]bool

var rPrincipleSet RPrincipleSet

// Add 方法用于将一个切片的元素添加到集合中
func (set RPrincipleSet) Add(items ...string) {
	for _, item := range items {
		set[item] = true
	}
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
	for item := range rPrincipleSet {
		fmt.Println(item)
	}
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
		var rPrincipalAssigned bool // 标记该时间段是否已经分配了前台负责人

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

		// 首先确保前台分配一个前台负责人
		if i%5 != 4 {
			for j := 0; j < len(availableAssistants); j++ {
				if _, isRPrincipal := rPrincipleSet[availableAssistants[j].Name]; isRPrincipal {
					timeSlotAssignments[i] = append(timeSlotAssignments[i], availableAssistants[j].Name)
					assistantAssignments[availableAssistants[j].Name]++
					rPrincipalAssigned = true
					break
				}
			}
		}

		if !rPrincipalAssigned {
			fmt.Println("警告：没有可用的负责人在当前的时段", i)

		}

		// 分配其它助理
		maxAssignments := 4
		if i%5 == 4 {
			// 对应 19:00-21:00 时间段
			maxAssignments = 3
		}
		for j := 0; j < len(availableAssistants) && len(timeSlotAssignments[i]) < maxAssignments; j++ {
			// 检查这个助理是否已被分配到这个时间段
			alreadyAssigned := false
			for _, assignedAssistant := range timeSlotAssignments[i] {
				if assignedAssistant == availableAssistants[j].Name {
					alreadyAssigned = true
					break
				}
			}
			// 如果这个助理没有被分配到这个时间段，则分配他/她
			if !alreadyAssigned {
				timeSlotAssignments[i] = append(timeSlotAssignments[i], availableAssistants[j].Name)
				assistantAssignments[availableAssistants[j].Name]++
			}
		}
	}
}

func assign() {
	preProcess()
	assignSDRPricipal()
	assignAssistant()
}
