package enums

import (
	"teagans-web-api/app/utilities"
	"strings"
)

type TaskPriority int

const (
	LOW		TaskPriority = iota
	MEDIUM
	HIGH
	URGENT
)

var priorityStrings = []string{"low", "medium", "high", "urgent"}

func(priority TaskPriority) String() string {
	return strings.Title(priorityStrings[priority])
}

func(priority TaskPriority) IsValid() bool {
	return int(priority) >= 0 && int(priority) < len(priorityStrings)
}

func ParseTaskPriorityString(priorityStr string) (TaskPriority, bool) {
	ndx := utilities.StringIndexOf(priorityStrings, strings.ToLower(priorityStr))

	if ndx >= 0 {
		return TaskPriority(ndx), true
	} else {
		return MEDIUM, false
	}
}

func ValToTaskPriority(val interface{}) (TaskPriority, bool) {
	var valInt		int
	var valFloat32 	float32
	var valFloat64 	float64
	var valStr		string
	var ok 			bool

	valInt, ok = val.(int)
	if ok {
		return TaskPriority(valInt), true
	}
	valFloat32, ok = val.(float32)
	if ok {
		return TaskPriority(valFloat32), true
	}
	valFloat64, ok = val.(float64)
	if ok {
		return TaskPriority(valFloat64), true
	}
	valStr, ok = val.(string)
	if ok {
		return ParseTaskPriorityString(valStr)
	}

	return MEDIUM, false
}
