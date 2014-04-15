package server

import (
	"github.com/campadrenalin/djdns/model"
	"github.com/miekg/dns"
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

func setupTestData() DjdnsServer {
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
	return s
}

func Test_DjdnsServer_GetRecords(t *testing.T) {
	// Setup
	s := setupTestData()

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

type ResolveTest struct {
	QuestionSection []dns.Question
	ExpectedAnswers []string
}

func (rt *ResolveTest) TestHandle(t *testing.T, s DjdnsServer) {
	// Construct query
	query := new(dns.Msg)
	query.Question = rt.QuestionSection

	// Get response
	response, err := s.Handle(query)
	if err != nil {
		t.Fatal(err)
	}

	// Construct expected response
	expected := new(dns.Msg)
	expected.Question = query.Question
	expected.Answer = make([]dns.RR, len(rt.ExpectedAnswers))
	for i, answer := range rt.ExpectedAnswers {
		rr, err := dns.NewRR(answer)
		if err != nil {
			t.Fatal(err)
		}
		expected.Answer[i] = rr
	}

	// Confirm equality
	if !reflect.DeepEqual(response, expected) {
		t.Log(response)
		t.Log(expected)
		t.Fatal("Response not equal to expected response")
	}
}

var resolve_tests = []ResolveTest{
	ResolveTest{
		QuestionSection: []dns.Question{
			dns.Question{
				"abcdef", dns.TypeA, dns.ClassINET},
		},
		ExpectedAnswers: []string{
			"first. A 1.1.1.1",
			"second. A 2.2.2.2",
		},
	},
	ResolveTest{
		QuestionSection: []dns.Question{
			dns.Question{
				"def", dns.TypeA, dns.ClassINET},
		},
		ExpectedAnswers: []string{},
	},
}

func Test_DjdnsServer_Handle(t *testing.T) {
	s := setupTestData()
	for _, test := range resolve_tests {
		test.TestHandle(t, s)
	}
}
