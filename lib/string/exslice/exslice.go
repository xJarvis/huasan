package exslice

func InStringArray(array *[]string,item string) bool {
	for _,value := range *array {
		if value == item {
			return true
		}
	}
	return false
}

func InStringArrayEx(array *[]string,items *[]string) bool {
	result := false
	for _, item := range *items {
		result = false
		for _, value := range *array {
			if value == item {
				result = true
			}
		}
		if !result {
			break
		}
	}
	return false
}