// main.go
package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

const (
	SLEEP_TIME          = 3 * time.Second
	KEYBOARD_BUFER_SIZE = 10000
	KEYMAP_BUFFER_SIZE  = 32000
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
	return make(map[uint8]string)
}

// Extract pressed keys from bufer buf
// It returns slice with key numbers in the same order
func GetKeyNumsFromOutput(buf []byte) []uint8 {
	const KEY_NUM_STRING_RE = "press[ ]+(\\d+)"
	re := regexp.MustCompile(KEY_NUM_STRING_RE)
	res_byte := re.FindAll(buf, -1)
	keyNums := make([]uint8, len(res_byte))
	re = regexp.MustCompile("\\d+")
	for i, line := range res_byte {
		num_byte := re.Find(line)
		if num, err := strconv.Atoi(string(num_byte)); err == nil {
			keyNums[i] = uint8(num)
		} else {
			log.Fatal(err)
		}
	}
	return keyNums
}

func main() {

	keyboard_id := 14
	cmd := exec.Command("xinput", "test", strconv.Itoa(keyboard_id))

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
