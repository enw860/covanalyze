package enricher

func ifNested(x, y int) bool {
	if x > 0 {
		if y > 0 {
			return true
		}
	}
	return false
}
