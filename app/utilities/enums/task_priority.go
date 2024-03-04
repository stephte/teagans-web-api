package enums

import (
	"teagans-web-api/app/utilities"
	"strings"
)

type TaskPriority int

const (
	LOW			TaskPriority = iota + 1
	MEDIUM
	HIGH
	URGENT
)

var priorityStrings = []string{"low", "medium", "high", "urgent"}

func(priority TaskPriority) String() string {
	if priority.IsValid() {
		return strings.Title(priorityStrings[priority-1])
	}

	return ""
}

func(priority TaskPriority) IsValid() bool {
	return int(priority) > 0 && int(priority) <= len(priorityStrings)
}

func NewTaskPriority(num int64) (TaskPriority, bool) {
	rv := TaskPriority(num)

	if rv.IsValid() {
		return rv, true
	} else {
		return LOW, false
	}
}

func ParseTaskPriorityString(str string) (TaskPriority, bool) {
	ndx := utilities.IndexOf(priorityStrings, strings.ToLower(str))

	if ndx >= 0 {
		return TaskPriority(ndx + 1), true
	} else {
		return LOW, false
	}
}
