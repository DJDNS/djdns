package server

import (
	"github.com/campadrenalin/djdns/model"
	"reflect"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer()
	if s.Port != 9953 {
		t.Fatalf("Expected port 9953, got %d", s.Port)
	}
	num_exp_branches := 0
	num_branches := len(s.Root.Branches)
	if num_branches != num_exp_branches {
		t.Fatalf(
			"Expected %d branches, got %d",
			num_exp_branches,
			num_branches,
		)
	}
}

type GetRecordsTest struct {
	Query    string
	Expected []model.Record
	ErrorMsg string
}

func (grt *GetRecordsTest) Run(t *testing.T, s DjdnsServer) {
	result := s.GetRecords(grt.Query)
	if !reflect.DeepEqual(result, grt.Expected) {
		t.Log(grt.Query)
		t.Log(grt.Expected)
		t.Log(result)
		t.Fatal(grt.ErrorMsg)
	}
}

func Test_DjdnsServer_GetRecords(t *testing.T) {
	// Setup
	s := NewServer()
	s.Root.Branches = []model.Branch{
		model.Branch{
			Selector: "abc",
			Records: []model.Record{
				model.Record{
					DomainName: "first",
					Rdata:      "1.1.1.1",
				},
				model.Record{
					DomainName: "second",
					Rdata:      "2.2.2.2",
				},
			},
		},
	}
	s.Root.Normalize()

	// Actual tests
	tests := []GetRecordsTest{
		GetRecordsTest{
			"abcde",
			s.Root.Branches[0].Records,
			"Basic request",
		},
		GetRecordsTest{
			"no such branch",
			nil,
			"Branch does not exist",
		},
	}
	for i := range tests {
		tests[i].Run(t, s)
	}
}
