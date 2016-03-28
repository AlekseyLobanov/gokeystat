// saving.go
package main

import (
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

// SaveToCsvWriter saves data to writerOut
// if fullExport saves log for each key else only sum for keys
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

// SaveToCsvFile saves data to path
// if fullExport saves log for each key else only sum for keys
func SaveToCsvFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	csvfile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer csvfile.Close()

	SaveToCsvWriter(data, keyMap, csvfile, fullExport)
}

// SaveToJSONWriter saves data to writerOut
// if fullExport saves log for each key else only sum for keys
// Save in one Json array
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

// SaveToJSONFile saves data to path
// if fullExport saves log for each key else only sum for keys
// Save in one Json array
func SaveToJSONFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	jsonFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	SaveToJSONWriter(data, keyMap, jsonFile, fullExport)
}

// SaveToJSONWriter saves data to writerOut
// if fullExport saves log for each key else only sum for keys
// Each record on new line
func SaveToJSLWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, fullExport bool) {
	type JSLStatForTime struct {
		Time int64
		Keys map[string]int
	}

	table := make([]JSLStatForTime, len(data))
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

	for ind, line := range table {
		lineBytes, err := json.Marshal(line)
		if err != nil {
			log.Fatal(err)
		}
		writerOut.Write(lineBytes)
		if ind != len(table)-1 {
			writerOut.Write([]byte("\n"))
		}
	}
}

// SaveToJSONWriter saves data to path
// if fullExport saves log for each key else only sum for keys
// Each record on new line
func SaveToJSLFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	jslFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jslFile.Close()

	SaveToJSLWriter(data, keyMap, jslFile, fullExport)
}

// SaveToCsvGzWriter same as SaveToCsvWriter but gunzip file before saving
func SaveToCsvGzWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, fullExport bool) {
	gzipWriter := gzip.NewWriter(writerOut)
	defer gzipWriter.Close()

	SaveToCsvWriter(data, keyMap, gzipWriter, fullExport)
}

// SaveToCsvGzFile same as SaveToCsvFile but gunzip file before saving
func SaveToCsvGzFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	jsonFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	gzipWriter := gzip.NewWriter(jsonFile)
	defer gzipWriter.Close()

	SaveToCsvWriter(data, keyMap, gzipWriter, fullExport)
}

// SaveToJSONGzWriter same as SaveToJSONWriter but gunzip file before saving
func SaveToJSONGzWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, fullExport bool) {
	gzipWriter := gzip.NewWriter(writerOut)
	defer gzipWriter.Close()

	SaveToJSONWriter(data, keyMap, gzipWriter, fullExport)
}

// SaveToJSONGzFile same as SaveToJSONFile but gunzip file before saving
func SaveToJSONGzFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	jsonFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	gzipWriter := gzip.NewWriter(jsonFile)
	defer gzipWriter.Close()

	SaveToJSONWriter(data, keyMap, gzipWriter, fullExport)
}

// SaveToJSLGzWriter same as SaveToJSLWriter but gunzip file before saving
func SaveToJSLGzWriter(data []StatForTime, keyMap map[uint8]string, writerOut io.Writer, fullExport bool) {
	gzipWriter := gzip.NewWriter(writerOut)
	defer gzipWriter.Close()

	SaveToJSLWriter(data, keyMap, gzipWriter, fullExport)
}

// SaveToJSLFile same as SaveToJSLFile but gunzip file before saving
func SaveToJSLGzFile(data []StatForTime, keyMap map[uint8]string, path string, fullExport bool) {
	jslFile, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer jslFile.Close()

	SaveToJSLGzWriter(data, keyMap, jslFile, fullExport)
}
