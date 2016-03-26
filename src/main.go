// main.go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SLEEP_TIME          = 3 * time.Second
	KEYBOARD_BUFER_SIZE = 10000
	DATABASE_NAME       = "./gokeystat.db"
	CAPTURE_TIME        = 5 * time.Second // time between capturing keyboard to db
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

func InitDb(db *sql.DB, keyMap map[uint8]string) {
	keyNums := make([]int, 0, len(keyMap))
	for keyNum := range keyMap {
		keyNums = append(keyNums, int(keyNum))
	}

	sqlInit := `CREATE TABLE IF NOT EXISTS keylog (
        time INTEGER primary key`

	for keyNum := range keyNums {
		sqlInit += ",\n" + "KEY" + strconv.Itoa(keyNum) + " INTEGER"
	}
	sqlInit += "\n);"
	_, err := db.Exec(sqlInit)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlInit)
	}

	// Inserting keymap to table
	sqlInit = `CREATE TABLE IF NOT EXISTS keymap (
        num INTEGER primary key,
		value STRING
	);`

	_, err = db.Exec(sqlInit)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlInit)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into keymap(num, value) values(?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for keyNum, keyName := range keyMap {
		_, err = stmt.Exec(keyNum, keyName)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

func main() {

	keyboardId := flag.Int("id", -1, "Your keyboard id")
	outputPath := flag.String("o", "", "Path to export file")
	flag.Parse()
	log.Println("keyboardId =", *keyboardId, "outputPath =", *outputPath)
	switch {
	case *keyboardId == -1 && *outputPath == "":
		flag.PrintDefaults()
		return
	case *keyboardId != -1:
		// Opening database
		db, err := sql.Open("sqlite3", DATABASE_NAME)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		keyMap := GetKeymap()

		InitDb(db, keyMap)
		cmd := exec.Command("xinput", "test", strconv.Itoa(*keyboardId))

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
	case *outputPath != "":
		//exporting here
	}
}
