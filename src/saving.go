// saving.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

func SaveToCsvWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, fullExport bool) {

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
	if fullExport {
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
			if fullExport {
				line = append(line, strconv.Itoa(rec.keys[key]))
			}
		}
		line = append(line, strconv.Itoa(sum))
		table = append(table, line)
	}

	writer := csv.NewWriter(writerOut)
	writer.WriteAll(table)
}

func SaveToCsvFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	csvfile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer csvfile.Close()

	SaveToCsvWriter(data, keyMap, csvfile, fullExport)
}

func SaveToJSONWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, fullExport bool) {
	type JSONStatForTime struct {
		Time int64
		Keys map[string]int
	}

	table := make([]JSONStatForTime, len(data))
	for i, stat := range data {
		table[i].Keys = make(map[string]int)

		table[i].Time = stat.time
		var sum int
		for numKey, key := range keyMap {
			if fullExport {
				table[i].Keys[key] = stat.keys[numKey]
			}
			sum += stat.keys[numKey]
		}
		table[i].Keys["sum"] = sum
	}

	outString, err := json.Marshal(table)
	if err != nil {
		log.Fatal(err)
	}

	writerOut.Write(outString)
}

func SaveToJSONFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	jsonFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	SaveToJSONWriter(data, keyMap, jsonFile, fullExport)
}
