package enums


type Enum interface {
	String()			string
	IsValid()			bool
}


func GetParseMethodsMap() map[string]interface{} {
	return map[string]interface{} {
		"TaskPriority": ParseTaskPriorityString,
		"TaskStatus": ParseTaskStatusString,
		"UserRole": ParseUserRoleString,
	}
}

func GetNewMethodsMap() map[string]interface{} {
	return map[string]interface{} {
		"TaskPriority": NewTaskPriority,
		"TaskStatus": NewTaskStatus,
		"UserRole": NewUserRole,
	}
}
