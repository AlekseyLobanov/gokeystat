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
	SLEEP_TIME          = 3 * time.Second // time between processing xinput output
	KEYBOARD_BUFER_SIZE = 10000
	DATABASE_NAME       = "file:gokeystat.db?cache=shared&mode=rwc"
	CAPTURE_TIME        = 5 // time in seconds between capturing keyboard to db
)

// StatForTime stotres pressed keys and beginning time
type StatForTime struct {
	time int64
	keys map[uint8]int
}

func (stat *StatForTime) Init() {
	stat.time = time.Now().Unix()
	stat.keys = make(map[uint8]int)
}

// GetKeymap returns map from key numbers to key names like "F1", "Tab", "d"
func GetKeymap() map[uint8]string {
	return GetKeymapFromOutput(GetKeymapOutput())
}

// GetKeymapOutput returns output of utility that prints system keymap
func GetKeymapOutput() []byte {
	cmd := exec.Command("xmodmap", "-pke")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

// GetKeymapFromOutput returns map with keymap from text
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

// GetKeyNumsFromOutput extract pressed keys from bufer buf
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

// InitDb creates tables, inserts keymap to db
func InitDb(db *sql.DB, keyMap map[uint8]string) {
	keyNums := GetKeyNumsFromKeyMap(keyMap)

	sqlInit := `CREATE TABLE IF NOT EXISTS keylog (
        time INTEGER primary key`

	for keyNum := range keyNums {
		sqlInit += ",\n" + "KEY" + strconv.Itoa(keyNum) + " INTEGER"
	}
	sqlInit += "\n);"

	// Inserting keymap to table
	sqlInit += ` CREATE TABLE IF NOT EXISTS keymap (
        num INTEGER primary key,
		value STRING
	);`

	_, err := db.Exec(sqlInit)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlInit)
	}

	rows, err := db.Query("SELECT COUNT(*) FROM keymap")
	if err != nil {
		log.Fatal(err)
	}

	var rowsCount int
	rows.Next()
	rows.Scan(&rowsCount)
	if rowsCount > 0 {
		// already inserted keymap
		return
	}
	rows.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO keymap(num, value) VALUES(?, ?)")
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
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	tx.Commit()
}

// GetStatTimesFromDb returns slice with StatForTime objects that
func GetStatTimesFromDb(db *sql.DB, fromTime int64, keyMap map[uint8]string) []StatForTime {
	sqlStmt := "select * from keylog where time > " + strconv.FormatInt(fromTime, 10)
	rows, err := db.Query(sqlStmt)
	if err != nil {
		log.Fatalln("Error with query", sqlStmt, " is ", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Fatalln("Failed to get columns", err)
	}

	rawResult := make([][]byte, len(cols))
	result := make([]int64, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	// keyNums[i] stores i-th keynum
	keyNums := GetKeyNumsFromKeyMap(keyMap)

	// result
	res := make([]StatForTime, 0)
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			log.Fatalln("Failed to scan row", err)
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = 0
			} else {
				// Only numbers in db: converting it to int64
				result[i], err = strconv.ParseInt(string(raw), 10, 64)
				if err != nil {
					log.Fatalln("Error when parsing ", raw, " from db:", err)
				}
			}
		}

		var resStatTime StatForTime
		resStatTime.time = result[0]
		resStatTime.keys = make(map[uint8]int)
		for index, val := range result[1:] {
			if val == 0 {
				continue
			}
			resStatTime.keys[uint8(keyNums[index])] = int(val)
		}
		res = append(res, resStatTime)
	}
	if err = rows.Err(); err != nil {
		log.Fatalln("Error when iterating over rows", err)
	}

	return res
}

func GetFileType(path string) string {

}

func main() {

	keyboardID := flag.Int("id", -1, "Your keyboard id")
	outputPath := flag.String("o", "", "Path to export file")
	fullExport := flag.Bool("full", false, "Export full stats")
	flag.Parse()

	log.Println("keyboardID =", *keyboardID, "outputPath =", *outputPath)

	// Opening database
	db, err := sql.Open("sqlite3", DATABASE_NAME)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)
	defer db.Close()

	keyMap := GetKeymap()

	InitDb(db, keyMap)

	switch {
	case *keyboardID == -1 && *outputPath == "":
		flag.PrintDefaults()
		return
	case *keyboardID != -1:

		cmd := exec.Command("xinput", "test", strconv.Itoa(*keyboardID))

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		// output of xinput command
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
		exportingData := GetStatTimesFromDb(db, 0, keyMap)
		filetype := GetFileType(*outputPath)
		log.Println(filetype)
		switch filetype {
		case ".csv":
			SaveToCsvFile(exportingData, keyMap, *outputPath, *fullExport)
		case ".json":
			SaveToJSONFile(exportingData, keyMap, *outputPath, *fullExport)
		case ".jsl":
			SaveToJSLFile(exportingData, keyMap, *outputPath, *fullExport)

		case ".csv.gz":
			SaveToCsvGzFile(exportingData, keyMap, *outputPath, *fullExport)
		case ".json.gz":
			SaveToJSONGzFile(exportingData, keyMap, *outputPath, *fullExport)
		case ".jsl.gz":
			SaveToJSLGzFile(exportingData, keyMap, *outputPath, *fullExport)
		default:
			log.Fatal("Incorrect file type")
		}
	}
}
