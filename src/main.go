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
)

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
