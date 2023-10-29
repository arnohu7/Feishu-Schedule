package main

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

// scheduleMap 用来存放助理提交的班表
var scheduleMap map[string][int(NumSlots) * int(NumDays)]bool

// timeSlotAssignments 用来存储每个时间段的助理分配情况
var timeSlotAssignments [int(NumSlots) * int(NumDays)][]string

// assistantAssignments 用来存储每个助理已经被分配的时间段数量
var assistantAssignments map[string]int
