package model

import "regexp"

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

// Test whether a branch matches a query.
//
// The branch's Selector property is used as a regex -
// if the query string passes that regex, the branch
// matches.
func (sb *Branch) Matches(query string) (bool, error) {
	return regexp.MatchString(sb.Selector, query)
}
