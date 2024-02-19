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
	return strings.Title(priorityStrings[priority-1])
}

func(priority TaskPriority) IsValid() bool {
	return int(priority) > 0 && int(priority) <= len(priorityStrings)
}

func ParseTaskPriorityString(priorityStr string) (TaskPriority, bool) {
	ndx := utilities.IndexOf(priorityStrings, strings.ToLower(priorityStr)) + 1

	if TaskPriority(ndx).IsValid() {
		return TaskPriority(ndx), true
	} else {
		return LOW, false
	}
}
