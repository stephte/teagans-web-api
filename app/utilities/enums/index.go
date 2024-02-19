package enums


type Enum interface {
	String()			string
	IsValid()			bool
}


func GetParseMethodsMap() map[string]interface{} {
	return map[string]interface{} {
		"ParseTaskPriorityString": ParseTaskPriorityString,
		"ParseTaskStatusString": ParseTaskStatusString,
		"ParseUserRoleString": ParseUserRoleString,
	}
}
