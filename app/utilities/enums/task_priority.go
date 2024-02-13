package enums

import (
	"strings"
)

type TaskPriority int

const (
	LOW		TaskPriority = iota + 1
	MEDIUM
	HIGH
	URGENT
)

var TaskPrioritiesArr []TaskPriority = []TaskPriority{LOW, MEDIUM, HIGH, URGENT}

var priorityStrings = []string{"low", "medium", "high", "urgent"}

func(priority TaskPriority) String() string {
	return priorityStrings[priority]
}

func(priority TaskPriority) IsValid() bool {
	return int(priority) <= len(priorityStrings) && int(priority) > 0
}

var roleMap = map[string]UserRole {
	"low":			LOW,
	"medium":		MEDIUM,
	"high":			HIGH,
	"urgent":		URGENT,
}

func ParseTaskPriority(roleStr string) (TaskPriority, bool) {
	priority, ok := priorityMap[strings.ToLower(priorityStr)]
	return priority, ok
}

func ValToPriority(val interface{}) (TaskPriority, bool) {
	var priorityInt int
	var priorityFloat32 float32
	var priorityFloat64 float64
	var priorityStr		string
	var ok bool

	priorityInt, ok = val.(int)
	if ok {
		return TaskPriority(priorityInt), true
	}
	priorityFloat32, ok = val.(float32)
	if ok {
		return TaskPriority(priorityFloat32), true
	}
	priorityFloat64, ok = val.(float64)
	if ok {
		return TaskPriority(priorityFloat64), true
	}
	priorityStr, ok = val.(string)
	if ok {
		return ParseTaskPriority(priorityStr)
	}

	return MEDIUM, false
}
