package utils

func DeduplicateInt64Slice(args ...[]int64) []int64 {
	numMap := map[int64]bool{}
	for _, slice := range args {
		for _, num := range slice {
			numMap[num] = true
		}
	}
	result := make([]int64, 0, len(numMap))
	for num := range numMap {
		result = append(result, num)
	}
	return result
}

// 求交集
func IntersectInt64Slice(args ...[]int64) []int64 {
	if len(args) == 0 {
		return []int64{}
	}
	numMap := map[int64]bool{}
	for _, num := range args[0] {
		numMap[num] = true
	}
	for i := 1; i < len(args); i++ {
		tempMap := map[int64]bool{}
		for _, num := range args[i] {
			if _, ok := numMap[num]; ok {
				tempMap[num] = true
			}
		}
		numMap = tempMap
	}
	result := make([]int64, 0, len(numMap))
	for num := range numMap {
		result = append(result, num)
	}
	return result
}
