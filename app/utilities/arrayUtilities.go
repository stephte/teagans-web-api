package utilities


func StringArrContains(arr []string, value string) bool {
	return StringIndexOf(arr, value) >= 0
}

func IntArrContains(arr []int64, value int64) bool {
	return IntIndexOf(arr, value) >= 0
}

func StringIndexOf(arr []string, value string) int64 {
	for ndx, el := range arr {
		if el == value {
			return int64(ndx)
		}
	}

	return -1
}

func IntIndexOf(arr []int64, value int64) int64 {
	for ndx, el := range arr {
		if el == value {
			return int64(ndx)
		}
	}

	return -1
}
