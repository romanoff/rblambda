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
}

func TestForceOldSyntax(t *testing.T) {
	for _, tt := range syntaxtests {
		out, changed := forceOldSyntax([]byte(tt.in))
		if string(out) != tt.out || changed != tt.changed {
			t.Errorf("Expected to get ['%v', %v], but got: ['%v', %v]", tt.out, tt.changed, string(out), changed)
		}
	}
}
