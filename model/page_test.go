package model

import (
	"reflect"
	"testing"
)

func testCompare(t *testing.T, msg string, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("a: %#v", a)
		t.Errorf("b: %#v", b)
		t.Fatalf(msg)
	}
}

func TestGetDemoData(t *testing.T) {
	data, err := GetJSONFromFile("demo.json")
	if err != nil {
		t.Fatal(err)
	}
	var sp Page
	err = sp.LoadFrom(data)
	if err != nil {
		t.Fatal(err)
	}

	expected_meta := PageMeta{
		Authority: "Campadrenalin/Rainfly_X/Philip Horger",
		Contact:   "@campadrenalin on github, or philip@roaming-initiative.com, or rainfly_x on reddit/HypeIRC/socialno.de",
		About:     "Roaming Initiative company DJDNS page.",
		Policy:    "Private",
	}
	testCompare(t, "Page meta differs", sp.Meta, expected_meta)

	// Could bulk-compare, but it makes it hard to track issues
	expected_branches := []Branch{
		Branch{
			Selector: "^(gmg\\.)?ri\\.hype\\.$",
			Targets:  []string{},
			Records: []Record{
				Record{
					DomainName: "ri.hype",
					Rtype:      "A",
					Rttl:       3600,
					Rdata:      "173.255.210.202",
				},
				Record{
					DomainName: "ri.hype",
					Rtype:      "AAAA",
					Rttl:       3600,
					Rdata:      "fcd5:7d07:2146:f18f:f937:d46e:77c9:80e7",
				},
				Record{
					DomainName: "ri.hype",
					Rtype:      "AAAA",
					Rttl:       3600,
					Rdata:      "2600:3c01::f03c:91ff:feae:1082",
				},
			},
		},
		Branch{
			Selector: "^orchard\\.ri\\.hype\\.$",
			Targets:  []string{},
			Records: []Record{
				Record{
					DomainName: "orchard.ri.hype",
					Rtype:      "A",
					Rttl:       3600,
					Rdata:      "106.186.18.242",
				},
				Record{
					DomainName: "orchard.ri.hype",
					Rtype:      "AAAA",
					Rttl:       3600,
					Rdata:      "fcd4:1dc1:cc08:c97d:85e2:6cad:eab8:864",
				},
			},
		},
	}
	testCompare(t, "Number of branches differs",
		len(sp.Branches),
		len(expected_branches),
	)
	for i, expected_branch := range expected_branches {
		branch := sp.Branches[i]
		testCompare(t, "Branches inequal", branch, expected_branch)
	}
}

func TestPage_LoadFrom_BadData(t *testing.T) {
	data := make(chan int)
	var sp Page
	err := sp.LoadFrom(data)
	if err == nil {
		t.Fatal("Page.LoadFrom should have failed")
	}
}

func TestPage_GetBranchForQuery_Multiple(t *testing.T) {
	var page Page
	page.Branches = []Branch{
		Branch{
			"xyz",
			make([]string, 0),
			make([]Record, 0),
		},
		Branch{
			"xy",
			make([]string, 0),
			make([]Record, 0),
		},
		Branch{
			"*?",
			make([]string, 0),
			make([]Record, 0),
		},
		Branch{
			"abc",
			make([]string, 0),
			make([]Record, 0),
		},
	}
	// Should match first, even though it could go for either
	// of the first two purely on a validity basis.
	branch := page.GetBranchForQuery("xyzf")
	expected := &page.Branches[0]
	if branch != expected {
		t.Fatalf("Expected %v, got %v", expected, branch)
	}
	// Can only match second - so it matches that.
	branch = page.GetBranchForQuery("xyf")
	expected = &page.Branches[1]
	if branch != expected {
		t.Fatalf("Expected %v, got %v", expected, branch)
	}
	// Matches the last one, even after a broken selector
	branch = page.GetBranchForQuery("abcde")
	expected = &page.Branches[3]
	if branch != expected {
		t.Fatalf("Expected %v, got %v", expected, branch)
	}
}

func TestPage_GetBranchForQuery_OneBranch(t *testing.T) {
	var page Page
	page.Branches = []Branch{
		Branch{
			"xyz",
			make([]string, 0),
			make([]Record, 0),
		},
	}
	branch := page.GetBranchForQuery("Some query")
	if branch != nil {
		t.Fatal("Returned *Branch should have been nil")
	}
	branch = page.GetBranchForQuery("xyz")
	if branch != &page.Branches[0] {
		t.Fatalf("Expected %v, got %v", &page.Branches[0], branch)
	}
}

func TestPage_GetBranchForQuery_NoBranches(t *testing.T) {
	var page Page
	branch := page.GetBranchForQuery("Some query")
	if branch != nil {
		t.Fatal("Returned *Branch should have been nil")
	}
}

func TestPage_GetBranchForQuery_Complex(t *testing.T) {
	data, err := GetJSONFromFile("demo.json")
	if err != nil {
		t.Fatal(err)
	}
	var sp Page
	err = sp.LoadFrom(data)
	if err != nil {
		t.Fatal(err)
	}

	confirm_correct_branch := func(qs []string, b *Branch) {
		for _, query := range qs {
			branch := sp.GetBranchForQuery(query)
			if branch != b {
				t.Errorf("Failed on query '%s'", query)
				t.Fatalf("Expected %v, got %v", b, branch)
			}
		}
	}

	confirm_correct_branch(
		[]string{
			"ri.hype.",
			"gmg.ri.hype.",
		},
		&sp.Branches[0],
	)
	confirm_correct_branch(
		[]string{
			"orchard.ri.hype.",
		},
		&sp.Branches[1],
	)
	confirm_correct_branch(
		[]string{
			"froot.loop.",
		},
		nil,
	)
}
