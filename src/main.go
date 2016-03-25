// main.go
package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"
)

const (
	SLEEP_TIME          = 3 * time.Second
	KEYBOARD_BUFER_SIZE = 10000
)

func GetKeymapFromOutput(buf []byte) map[uint8]string {
	return make(map[uint8]string)
}

// Extract pressed keys from bufer buf
// It returns slice with key numbers in the same order
func GetKeyNumsFromOutput(buf []byte) []uint8 {
	return make([]uint8, 0)
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
