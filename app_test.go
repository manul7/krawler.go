package main

import "testing"


func TestBuildBaseURL(t *testing.T) {
	var variants = []struct {
        url			string
        expected	string
		}{
        {"example.com", "http://example.com"}, // No scheme
		{"http://www.example.com", "http://www.example.com"}, // Standrad with www
		{"http://example.com", "http://example.com"}, // Standard
		{"example.com/fromthis", "http://example.com/fromthis"}, // URL with path
	}

	for _, u := range variants {
		r, err := buildBaseUrl(&u.url)
		if err != nil {
			t.Error("Error while building Base URL")
		}
		if r.String() != u.expected {
			t.Errorf("Error: %v != %v", r.String(), u.expected)
		}
	}
	
}

func TestBuildBaseURL_incorrect(t *testing.T) {

}
