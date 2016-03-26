// main.go
package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	SLEEP_TIME          = 3 * time.Second
	KEYBOARD_BUFER_SIZE = 10000
)

// Return map from key numbers to key names like "F1", "Tab", "d"
func GetKeymap() map[uint8]string {
	return GetKeymapFromOutput(GetKeymapOutput())
}

// Return output of utility that prints system keymap
func GetKeymapOutput() []byte {
	cmd := exec.Command("xmodmap", "-pke")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

// Return map with keymap from text
func GetKeymapFromOutput(buf []byte) map[uint8]string {
	const KEY_NUM_STRING_RE = "\\d+[ ]*=[ ]*\\S+"
	re := regexp.MustCompile(KEY_NUM_STRING_RE)
	resByte := re.FindAll(buf, -1)
	keyMap := make(map[uint8]string)
	for _, line := range resByte {
		lineSpitted := strings.Split(string(line), " ")
		if key, err := strconv.Atoi(lineSpitted[0]); err == nil {
			keyMap[uint8(key)] = lineSpitted[2]
		}
	}
	return keyMap
}

// Extract pressed keys from bufer buf
// It returns slice with key numbers in the same order
func GetKeyNumsFromOutput(buf []byte) []uint8 {
	const KEY_NUM_STRING_RE = "press[ ]+(\\d+)"
	re := regexp.MustCompile(KEY_NUM_STRING_RE)
	resByte := re.FindAll(buf, -1)
	keyNums := make([]uint8, len(resByte))
	re = regexp.MustCompile("\\d+")
	for i, line := range resByte {
		numByte := re.Find(line)
		if num, err := strconv.Atoi(string(numByte)); err == nil {
			keyNums[i] = uint8(num)
		} else {
			log.Fatal(err)
		}
	}
	return keyNums
}

func main() {

	keyboardID := 14
	cmd := exec.Command("xinput", "test", strconv.Itoa(keyboardID))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, KEYBOARD_BUFER_SIZE)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		// processing buf here
		fmt.Println(n)
		time.Sleep(SLEEP_TIME)
	}
}
