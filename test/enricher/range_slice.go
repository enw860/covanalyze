package enricher

func rangeSlice(items []int) int {
	sum := 0
	for i, v := range items {
		sum += i + v
	}
	return sum
}
