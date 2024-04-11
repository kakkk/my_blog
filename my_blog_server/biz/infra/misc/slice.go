package misc

import "my_blog/biz/infra/config"

// 去重
func SliceDeduplicate[T comparable](arr []T) []T {
	temp := make(map[T]struct{})
	l := len(arr)
	if l == 0 {
		return arr
	}
	res := make([]T, 0, l)
	for _, item := range arr {
		_, ok := temp[item]
		if ok {
			continue
		}
		temp[item] = struct{}{}
		res = append(res, item)
	}

	return res[:len(temp)]
}

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

func MapToSet[K comparable, T any](m map[K]T) []T {
	result := make([]T, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func MapToListByOrder[K comparable, T any](order []K, m map[K]T) []T {
	result := make([]T, 0, len(m))
	for _, k := range order {
		if v, ok := m[k]; ok {
			result = append(result, v)
		}
	}
	return result
}

func GetIDsByPage(ids []int64, page int64) ([]int64, bool) {
	if page < 1 {
		return []int64{}, false
	}
	if page == 1 {
		if len(ids) < config.GetPageListSize() {
			return ids, false
		}
		hasMore := true
		if len(ids) <= config.GetPageListSize() {
			hasMore = false
		}
		ids = ids[0:config.GetPageListSize()]
		result := make([]int64, len(ids))
		copy(result, ids)

		return result, hasMore
	}
	begin := (page - 1) * config.GetPageListSizeI64()
	end := begin + config.GetPageListSizeI64()
	hasMore := true
	if len(ids) <= int(end) {
		hasMore = false
	}
	if len(ids) < int(begin) {
		return []int64{}, false
	}
	if len(ids) < int(end) {
		ids = ids[begin:]
		result := make([]int64, len(ids))
		copy(result, ids)
		return result, false
	}
	ids = ids[begin:end]
	result := make([]int64, len(ids))
	copy(result, ids)
	return result, hasMore
}
