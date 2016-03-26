// saving.go
package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

func SaveToCsvWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, isOnlySum bool) {

	numKeysInt := make([]int, 0)
	for key := range keyMap {
		numKeysInt = append(numKeysInt, int(key))
	}
	sort.Ints(numKeysInt)
	numKeys := make([]uint8, 0)
	for _, key := range numKeysInt {
		numKeys = append(numKeys, uint8(key))
	}

	titleLine := make([]string, 0)
	titleLine = append(titleLine, "Time")
	if !isOnlySum {
		for _, key := range numKeys {
			titleLine = append(titleLine, keyMap[key])
		}
	}
	titleLine = append(titleLine, "Sum")

	table := make([][]string, 0)
	table = append(table, titleLine)
	for _, rec := range data {
		line := make([]string, 0)
		line = append(line, strconv.Itoa(int(rec.time)))
		var sum int
		for _, key := range numKeys {
			sum += rec.keys[key]
			if !isOnlySum {
				line = append(line, strconv.Itoa(rec.keys[key]))
			}
		}
		line = append(line, strconv.Itoa(sum))
		table = append(table, line)
	}

	writer := csv.NewWriter(writerOut)
	writer.WriteAll(table)
}

func SaveToCsvFile(data []StatForTime, keyMap map[uint8]string, path string, isOnlySum bool) {
	csvfile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer csvfile.Close()

	SaveToCsvWriter(data, keyMap, csvfile, isOnlySum)
}
