package main

import "testing"

func TestUnpackEmpty(t *testing.T) {
	output, err := unpack("")
	if output != "" || err != nil {
		t.Fatalf(`unpack("") = %q, %v, want "", nil`, output, err)
	}
}

func TestUnpackFirstDigit(t *testing.T) {
	output, err := unpack("4a4b5n2m1")
	if err == nil {
		t.Fatalf(`unpack with first digit gives no error and outputs %q`, output)
	}
}

func TestUnpack(t *testing.T) {
	test := `qwe\\\\\35\5\\3`
	want := `qwe\\333335\\\`
	output, err := unpack(test)
	if output != want || err != nil {
		t.Fatalf(`unpack(%q) = %q, %v, want %q, nil`, test, output, err, want)
	}
}
