package main

import "testing"

func TestRegExp(t *testing.T) {
	if !PatternLogFile.MatchString(`hello.Log`) {
		t.Fatal()
	}
	if !PatternCompressedLogFile.MatchString(`hello.log-22222.gz`) {
		t.Fatal()
	}
	if !PatternCompressedLogFile.MatchString(`hello.log22222.gz`) {
		t.Fatal()
	}
	if !PatternCompressedLogFile.MatchString(`hello.log22222.gz222`) {
		t.Fatal()
	}
	if !PatternCompressedLogFile.MatchString(`hello-2222.log22222.gz`) {
		t.Fatal()
	}
}
