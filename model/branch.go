package model

type Branch struct {
	Selector string
	Targets  []string
	Records  []Record
}

// Fill in data that was left implied via defaults
func (sb *Branch) Normalize() {
	for r := range sb.Records {
		sb.Records[r].Normalize()
	}
}
