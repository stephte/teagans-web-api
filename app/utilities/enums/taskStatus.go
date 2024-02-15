package enums

import (
	"teagans-web-api/app/utilities"
	"strings"
)

type TaskStatus int

const (
	TODO		TaskStatus = iota
	WAITING
	STARTED
	COMPLETE
)

var statusStrings = []string{"todo", "waiting", "started", "complete"}

func(status TaskStatus) String() string {
	return strings.Title(statusStrings[status])
}

func(status TaskStatus) IsValid() bool {
	return int(status) >= 0 && int(status) < len(statusStrings)
}

func ParseTaskStatusString(statusStr string) (TaskStatus, bool) {
	ndx := utilities.StringIndexOf(statusStrings, strings.ToLower(statusStr))

	if ndx >= 0 {
		return TaskStatus(ndx), true
	} else {
		return TODO, false
	}
}

func ValToTaskStatus(val interface{}) (TaskStatus, bool) {
	var valInt		int
	var valFloat32 	float32
	var valFloat64 	float64
	var valStr		string
	var ok 			bool

	valInt, ok = val.(int)
	if ok {
		return TaskStatus(valInt), true
	}
	valFloat32, ok = val.(float32)
	if ok {
		return TaskStatus(valFloat32), true
	}
	valFloat64, ok = val.(float64)
	if ok {
		return TaskStatus(valFloat64), true
	}
	valStr, ok = val.(string)
	if ok {
		return ParseTaskStatusString(valStr)
	}

	return TODO, false
}
