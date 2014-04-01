package main

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
	var sp SerializedPage
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
	expected_branches := []SerializedBranch{
		SerializedBranch{
			Selector: "^(gmg\\.)?ri\\.hype$",
			Targets:  []string{},
			Records: []SerializedRecord{
				SerializedRecord{
					DomainName: "ri.hype",
					Rtype:      "A",
					Rdata:      "173.255.210.202",
				},
				SerializedRecord{
					DomainName: "ri.hype",
					Rtype:      "AAAA",
					Rdata:      "fcd5:7d07:2146:f18f:f937:d46e:77c9:80e7",
				},
				SerializedRecord{
					DomainName: "ri.hype",
					Rtype:      "AAAA",
					Rdata:      "2600:3c01::f03c:91ff:feae:1082",
				},
			},
		},
		SerializedBranch{
			Selector: "^orchard\\.ri\\.hype$",
			Targets:  []string{},
			Records: []SerializedRecord{
				SerializedRecord{
					DomainName: "orchard.ri.hype",
					Rtype:      "A",
					Rdata:      "106.186.18.242",
				},
				SerializedRecord{
					DomainName: "orchard.ri.hype",
					Rtype:      "AAAA",
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
