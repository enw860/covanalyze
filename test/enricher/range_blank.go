package enricher

func rangeBlank(items []int) int {
	sum := 0
	for _, v := range items {
		sum += v
	}
	return sum
}
