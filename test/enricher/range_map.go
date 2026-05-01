package enricher

func rangeMap(myMap map[string]int) int {
	sum := 0
	for k, v := range myMap {
		sum += len(k) + v
	}
	return sum
}
