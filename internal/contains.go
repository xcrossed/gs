package internal

// ContainsString 在一个 string 数组中进行查找，找不到返回 -1。
func ContainsString(array []string, val string) int {
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			return i
		}
	}
	return -1
}
