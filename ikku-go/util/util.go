package util

func SContains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func IContains(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func Sum(a []int) int {
	var result int
	for _, n := range a {
		result += n
	}
	return result
}
