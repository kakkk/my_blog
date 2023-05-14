package cachex

func sliceDeduplicate[T comparable](arr []T) []T {
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
