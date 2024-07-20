package dsl

import (
	"os"
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	str, err := os.ReadFile("../scenes/test_scene")
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(string(str))
	if got, err := p.Parse(); err != nil {
		t.Errorf("NewParser() = %v, want %v", got, str)
	} else {
		if !reflect.DeepEqual(got, str) {
			t.Errorf("NewParser() = %v, want %v", got, str)
		}
	}
}
