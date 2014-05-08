package server

import "testing"

func TestFilePageGetter_GetPage(t *testing.T) {
	pg := NewFilePageGetter()
	filename := "../model/demo.json"
	page, err := pg.GetPage(filename, make(chan interface{}))
	if err != nil {
		t.Fatal(err)
	}
	if page.Url != filename {
		t.Fatal("Wrong URL property")
	}
	if page.Expired {
		t.Fatal("Page should not start out expired")
	}
	// Just use this as a sanity check, no thorough test
	expected_policy := "Private"
	got_policy := page.Data.Meta.Policy
	if got_policy != expected_policy {
		t.Fatalf("Expected %s, got %s", expected_policy, got_policy)
	}
}

func TestFilePageGetter_GetPage_NoSuchFile(t *testing.T) {
	pg := NewFilePageGetter()
	filename := "../model/nosuchfile.json"
	_, err := pg.GetPage(filename, make(chan interface{}))
	if err == nil {
		t.Fatal("Should have announced failure")
	}
}
