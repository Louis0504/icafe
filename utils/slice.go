package utils

func SetDiff(a, b []int64) []int64 {
	// Set Difference: A - B
	diff := []int64{}
	mapA := make(map[int64]bool)
	mapB := make(map[int64]bool)

	for _, item := range a {
		mapA[item] = true
	}

	for _, item := range b {
		mapB[item] = true
	}

	for key := range mapA {
		if _, ok := mapB[key]; !ok {
			diff = append(diff, key)
		}
	}
	return diff
}

func SetDiffStr(a, b []string) []string {
	// Set Difference: A - B
	diff := []string{}
	mapA := make(map[string]bool)
	mapB := make(map[string]bool)

	for _, item := range a {
		mapA[item] = true
	}

	for _, item := range b {
		mapB[item] = true
	}

	for key := range mapA {
		if _, ok := mapB[key]; !ok {
			diff = append(diff, key)
		}
	}
	return diff
}

func Int64SliceToMap(ids []int64) map[int64]bool {
	resMap := make(map[int64]bool)
	for _, value := range ids {
		resMap[value] = true
	}
	return resMap
}

func Int64RemoveDuplicate(slc []int64) []int64 {
	result := make([]int64, 0)
	tempMap := make(map[int64]int)
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		// 加入 map 后，map 长度变化，则元素不重复
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func ContainStr(slc []string, target string) bool {
	for _, element := range slc {
		if target == element {
			return true
		}
	}
	return false
}

func IndexOfInt64(slc []int64, target int64) int64 {
	for i, element := range slc {
		if element == target {
			return int64(i)
		}
	}
	return -1
}

func ContainInt64(slc []int64, target int64) bool {
	for _, element := range slc {
		if target == element {
			return true
		}
	}
	return false
}

func IndexSliceStr(slc []string, target string) int64 {
	for i, element := range slc {
		if target == element {
			return int64(i)
		}
	}
	return -1
}

func IndexSliceInt64(slc []int64, target int64) int64 {
	for i, element := range slc {
		if target == element {
			return int64(i)
		}
	}
	return -1
}

func RemoveSliceInt64(slc []int64, target int64) []int64 {
	index := IndexSliceInt64(slc, target)
	if index == -1 {
		return slc
	}
	return append(slc[:index], slc[index+1:]...)
}

func RemoveSliceStr(slc []string, target string) []string {
	index := IndexSliceStr(slc, target)
	if index == -1 {
		return slc
	}
	return append(slc[:index], slc[index+1:]...)
}

func Int64SliceToChunks(slc []int64, chunkSize int) [][]int64 {
	var chunks [][]int64
	if chunkSize <= 0 {
		panic("invaid chunkSize")
	}
	if len(slc) <= chunkSize {
		chunks = append(chunks, slc)
		return chunks
	}
	for i := 0; i < len(slc); i += chunkSize {
		end := i + chunkSize
		if end > len(slc) {
			end = len(slc)
		}
		chunks = append(chunks, slc[i:end])
	}
	return chunks
}

func SafeSliceCut(sli []int64, start, end int64) []int64 {
	start = MaxInt64(0, start)
	if end <= 0 || start > end {
		return []int64{}
	}

	sliLen := int64(len(sli))

	if start > sliLen {
		return []int64{}
	}

	if end > sliLen {
		return sli[start:sliLen]
	}

	return sli[start:end]
}

func RepeatInt64(n int64, c int) []int64 {
	v := make([]int64, c)
	for i := range v {
		v[i] = n
	}

	return v
}
