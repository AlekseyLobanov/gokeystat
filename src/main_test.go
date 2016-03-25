// main_test.go
package main

import (
	"reflect"
	"testing"
)

func TestGetKeymapFromOutput(t *testing.T) {
	var buf []byte
	var keymap map[uint8]string

	const test1 = "keycode  19 = 0 parenright 0 parenright\n" +
		"keycode  20 = minus underscore minus underscore\n" +
		"keycode  21 = equal plus equal plus"
	result1 := map[uint8]string{19: "0", 20: "minus", 21: "equal"}

	const test2 = "keycode 119 = Delete NoSymbol Delete\n" +
		"keycode 120 =\n" +
		"keycode 121 = XF86AudioMute NoSymbol XF86AudioMute"
	result2 := map[uint8]string{119: "Delete", 121: "XF86AudioMute"}

	// Test1. Simple
	buf = []byte(test1)
	keymap = GetKeymapFromOutput(buf)
	if !reflect.DeepEqual(keymap, result1) {
		t.Fail()
	}

	// Test2. With empty keys
	buf = []byte(test2)
	keymap = GetKeymapFromOutput(buf)
	if !reflect.DeepEqual(keymap, result2) {
		t.Fail()
	}
}

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
