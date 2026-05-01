package enricher

func forInfinite() int {
	count := 0
	for {
		count++
		if count > 5 {
			break
		}
	}
	return count
}
