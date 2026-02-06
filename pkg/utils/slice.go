package utils

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func RemoveString(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func UniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func FilterStrings(slice []string, predicate func(string) bool) []string {
	var result []string
	for _, s := range slice {
		if predicate(s) {
			result = append(result, s)
		}
	}
	return result
}

func MapStrings(slice []string, mapper func(string) string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = mapper(s)
	}
	return result
}

func ReduceStrings(slice []string, reducer func(string, string) string, initial string) string {
	result := initial
	for _, s := range slice {
		result = reducer(result, s)
	}
	return result
}

func ChunkStrings(slice []string, size int) [][]string {
	if size <= 0 {
		return nil
	}
	chunks := make([][]string, 0, (len(slice)+size-1)/size)
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func ReverseStrings(slice []string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[len(slice)-1-i] = s
	}
	return result
}

func MergeIntSlices(slices ...[]int) []int {
	length := 0
	for _, slice := range slices {
		length += len(slice)
	}
	result := make([]int, 0, length)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}

func MergeStringSlices(slices ...[]string) []string {
	length := 0
	for _, slice := range slices {
		length += len(slice)
	}
	result := make([]string, 0, length)
	for _, slice := range slices {
		result = append(result, slice...)
	}
	return result
}
