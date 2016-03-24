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
		fmt.Println(len(buf), n)
		time.Sleep(5 * time.Second)
	}
}
