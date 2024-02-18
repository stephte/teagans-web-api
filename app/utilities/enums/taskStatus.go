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
	return strings.Title(statusStrings[status])
}

func(status TaskStatus) IsValid() bool {
	return int(status) > 0 && int(status) <= len(statusStrings)
}

func ParseTaskStatusString(statusStr string) (TaskStatus, bool) {
	ndx := utilities.StringIndexOf(statusStrings, strings.ToLower(statusStr)) + 1

	if ndx >= 0 {
		return TaskStatus(ndx), true
	} else {
		return TODO, false
	}
}
