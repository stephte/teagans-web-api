package enums

import (
	"teagans-web-api/app/utilities"
	"strings"
)

type TaskStatus int

const (
	TODO		TaskStatus = iota + 1
	WAITING
	STARTED
	COMPLETE
)

var statusStrings = []string{"todo", "waiting", "started", "complete"}

func(status TaskStatus) String() string {
	if status.IsValid() {
		return strings.Title(statusStrings[status-1])
	}

	return ""
}

func(status TaskStatus) IsValid() bool {
	return int(status) > 0 && int(status) <= len(statusStrings)
}

func NewTaskStatus(num int64) (TaskStatus, bool) {
	rv := TaskStatus(num)

	if rv.IsValid() {
		return rv, true
	} else {
		return TODO, false
	}
}

func ParseTaskStatusString(str string) (TaskStatus, bool) {
	ndx := utilities.IndexOf(statusStrings, strings.ToLower(str))

	if ndx >= 0 {
		return TaskStatus(ndx + 1), true
	} else {
		return TODO, false
	}
}
