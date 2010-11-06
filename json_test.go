package main

import (
	"testing"
	"json"
)

func TestMarshall(t *testing.T) {
	bin := make(Bin, 10)

	_, ok := bin.(json.Marshaler)
	if !ok {
		t.Fatal("main.Bin not a json.Marshaler")
	}

	_, ok := bin.(json.Unmarshaler)
	if !ok {
		t.Fatal("main.Bin not a json.Unmarshaler")
	}
}
