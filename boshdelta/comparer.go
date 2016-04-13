package boshdelta

// Comparer compares two BOSH releases
type Comparer struct {
	Release1 string
	Release2 string
}

// Delta is the result from comparing two BOSH releases
type Delta struct {
}

// Compare two BOSH releases
func Compare(release1, release2 string) (*Delta, error) {
	return nil, nil
}
