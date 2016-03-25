// main_test.go
package main

import (
	"reflect"
	"testing"
)

func TestGetKeyNumsFromOutput(t *testing.T) {
	var buf []byte
	var keyNums []uint8

	const test1 = "key press 36\nkey release 41\nkey press   41"
	var result1 = []uint8{36, 41}

	const test2 = ""
	var result2 = []uint8{}

	// Test1. Simple
	buf = []byte(test1)
	keyNums = GetKeyNumsFromOutput(buf)
	if !reflect.DeepEqual(keyNums, result1) {
		t.Fail()
	}

	// Test2. Clear
	buf = []byte(test2)
	keyNums = GetKeyNumsFromOutput(buf)
	if !reflect.DeepEqual(keyNums, result2) {
		t.Fail()
	}
}
