package server

import (
	"fmt"
	"github.com/campadrenalin/djdns/model"
	"github.com/miekg/dns"
	"reflect"
	"testing"
	"time"
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

type ResolveFunc func(*dns.Msg) (*dns.Msg, error)

func (rt *ResolveTest) run(t *testing.T, resolver ResolveFunc) {
	// Construct query
	query := new(dns.Msg)
	query.Question = rt.QuestionSection

	// Get response
	response, err := resolver(query)
	if err != nil {
		t.Log(rt)
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
	expected.Ns = make([]dns.RR, 0)
	expected.Extra = make([]dns.RR, 0)

	compare_part := func(item1, item2 interface{}, name string) {
		if !reflect.DeepEqual(item1, item2) {
			t.Logf(" * %s is different", name)
			t.Logf("%#v vs %#v", item1, item2)
		}
	}

	// DNS package tends to be loose about some encoding details,
	// only calculating them right before putting the data on the
	// wire.
	sanitize := func(rr_list []dns.RR) {
		for i := range rr_list {
			rr_list[i].Header().Rdlength = 0
		}
	}
	for _, msg := range []*dns.Msg{response, expected} {
		sanitize(msg.Answer)
		sanitize(msg.Ns)
		sanitize(msg.Extra)
	}

	// Confirm equality
	if !reflect.DeepEqual(response, expected) {
		t.Log(response)
		t.Log(expected)
		t.Log("Response not equal to expected response")

		// More DRY to use reflect, but it would also be like
		// chewing broken glass.
		compare_part(response.MsgHdr, expected.MsgHdr, "MsgHdr")
		compare_part(response.Compress, expected.Compress, "Compress")
		compare_part(response.Question, expected.Question, "Question")
		compare_part(response.Answer, expected.Answer, "Answer")
		compare_part(response.Ns, expected.Ns, "Ns")
		compare_part(response.Extra, expected.Extra, "Extra")

		t.Log(response.Answer[0].Header().Rdlength)
		t.Log(expected.Answer[0].Header().Rdlength)
		t.FailNow()
	}
}

func (rt *ResolveTest) TestHandle(t *testing.T, s DjdnsServer) {
	rt.run(t, func(query *dns.Msg) (*dns.Msg, error) {
		return s.Handle(query)
	})
}

func (rt *ResolveTest) TestResolve(t *testing.T, c *dns.Client, addr string) {
	rt.run(t, func(query *dns.Msg) (*dns.Msg, error) {
		response, _, err := c.Exchange(query, addr)
		return response, err
	})
}

// TODO: Unqualified domains/faliure
var resolve_tests = []ResolveTest{
	ResolveTest{
		QuestionSection: []dns.Question{
			dns.Question{
				"abcdef.", dns.TypeA, dns.ClassINET},
		},
		ExpectedAnswers: []string{
			"first. A 1.1.1.1",
			"second. A 2.2.2.2",
		},
	},
	ResolveTest{
		QuestionSection: []dns.Question{
			dns.Question{
				"def.", dns.TypeA, dns.ClassINET},
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

func Test_DjdnsServer_Run(t *testing.T) {
	s := setupTestData()
	host, port := "127.0.0.1", 9953
	addr := fmt.Sprintf("%s:%d", host, port)

	go func() {
		t.Fatal(s.Run(addr))
	}()
	defer s.Close()
	<-time.After(50 * time.Millisecond)

	c := new(dns.Client)
	for _, test := range resolve_tests {
		test.TestResolve(t, c, addr)
	}
}
