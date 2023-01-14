package utilities


func StringArrContains(arr []string, value string) bool {
	for _, el := range arr {
		if el == value {
			return true
		}
	}

	return false
}

func IntArrContains(arr []int, value int) bool {
	for _, el := range arr {
		if el == value {
			return true
		}
	}

	return false
}
