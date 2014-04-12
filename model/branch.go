package model

type Branch struct {
	Selector string
	Targets  []string
	Records  []Record
}

// Set default values for missing data
func (sb *Branch) Normalize() {
	for r := range sb.Records {
		sb.Records[r].Normalize()
	}
}
