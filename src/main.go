// main.go
package main

import (
	"database/sql"
	"flag"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	SLEEP_TIME          = 3 * time.Second
	KEYBOARD_BUFER_SIZE = 10000
	DATABASE_NAME       = "./gokeystat.db"
	CAPTURE_TIME        = 5 // time in seconds between capturing keyboard to db
)

type StatForTime struct {
	time int64
	keys map[uint8]int
}

func (stat *StatForTime) Init() {
	stat.time = time.Now().Unix()
	stat.keys = make(map[uint8]int)
}

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

func GetKeyNumsFromKeyMap(keyMap map[uint8]string) []int {
	res := make([]int, 0, len(keyMap))
	for keyNum := range keyMap {
		res = append(res, int(keyNum))
	}
	sort.Ints(res)
	return res
}

func InitDb(db *sql.DB, keyMap map[uint8]string) {
	keyNums := GetKeyNumsFromKeyMap(keyMap)

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
	rows, err := db.Query("select COUNT(*) from keymap")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var rowsCount int
	rows.Next()
	rows.Scan(&rowsCount)
	if rowsCount > 0 {
		// already inserted keymap
		return
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

func AddStatTimeToDb(db *sql.DB, statTime StatForTime, keyMap map[uint8]string) {
	keyNums := GetKeyNumsFromKeyMap(keyMap)
	sqlStmt := "insert into keylog(time"
	for keyNum := range keyNums {
		sqlStmt += ",\n" + "KEY" + strconv.Itoa(keyNum)
	}
	sqlStmt += ") values "
	sqlStmt += "(" + strconv.FormatInt(statTime.time, 10)
	for keyNum := range keyNums {
		keyNumber, _ := statTime.keys[uint8(keyNum)]
		sqlStmt += ",\n" + strconv.Itoa(keyNumber)
	}
	sqlStmt += ")"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
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
		var curStat StatForTime
		curStat.Init()
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			// processing buf here
			for _, keyNum := range GetKeyNumsFromOutput(buf[:n]) {
				oldKeyCount, _ := curStat.keys[keyNum]
				curStat.keys[keyNum] = oldKeyCount + 1
			}

			// Every CAPTURE_TIME seconds save to BD
			if time.Now().Unix()-curStat.time > CAPTURE_TIME {
				AddStatTimeToDb(db, curStat, keyMap)
				curStat.Init()
			}

			time.Sleep(SLEEP_TIME)
		}
	case *outputPath != "":
		//exporting here
	}
}
