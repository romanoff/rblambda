package main

import (
	"testing"
)

var syntaxtests = []struct {
	in      string
	out     string
	changed bool
}{
	{"-> { 4 }", "lambda { 4 }", true},
	{"lambda { 4 }", "lambda { 4 }", false},
	{"-> (x) { 4 }", "lambda { |x| 4 }", true},
	{"-> () { 4 }", "lambda { 4 }", true},
	{`-> do |x|
a = 5 if x == true
end`, `lambda do |x|
a = 5 if x == true
end`, true},
	{`-> do |x|
  true if true
  if x == true
    a = 5
  end
end`, `lambda do |x|
  true if true
  if x == true
    a = 5
  end
end`, true},
}

func TestForceOldSyntax(t *testing.T) {
	for _, tt := range syntaxtests {
		out, changed := forceOldSyntax([]byte(tt.in))
		if string(out) != tt.out || changed != tt.changed {
			t.Errorf("Expected to get ['%v', %v], but got: ['%v', %v]", tt.out, tt.changed, string(out), changed)
		}
	}
}
