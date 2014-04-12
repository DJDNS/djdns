package model

import "testing"

type MatchTest struct {
	BranchSelector string
	Query          string
	ShouldSucceed  bool
	ShouldError    bool
}

func (mt *MatchTest) Run(t *testing.T) {
	var branch Branch
	branch.Selector = mt.BranchSelector

	match, err := branch.Matches(mt.Query)
	if !mt.ShouldError && err != nil {
		t.Log(mt.BranchSelector)
		t.Error(err)
	} else if mt.ShouldError && err == nil {
		t.Log(mt.BranchSelector)
		t.Error("err should not have been nil")
	}
	if match != mt.ShouldSucceed {
		t.Errorf(
			"%s => %s != %t",
			mt.BranchSelector,
			mt.Query,
			mt.ShouldSucceed,
		)
	}
}

func TestBranch_Matches(t *testing.T) {
	tests := []MatchTest{
		MatchTest{"ell", "hello", true, false},
		MatchTest{"hello", "el", false, false},
		MatchTest{"*?", "el", false, true},
		MatchTest{"example\\.com", "example.com", true, false},
		MatchTest{"example\\.com", "example0com", false, false},
		MatchTest{"example.com", "example0com", true, false},
	}
	for _, test := range tests {
		test.Run(t)
	}
}
