// saving_test.go
package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func GenerateRandStatsForTime(N int) []StatForTime {
	data := make([]StatForTime, N)
	keyMap := GetKeymap()
	rnd := rand.New(rand.NewSource(42))
	for i := 0; i < N; i++ {
		var curStat StatForTime
		curStat.Init()
		for keyNum := range keyMap {
			if rnd.Float32() > 0.7 {
				curStat.keys[keyNum] = rnd.Intn(5000)
			}
		}
		data = append(data, curStat)
	}

	return data
}

func BenchmarkCsvSavingOnlySum(b *testing.B) {
	data := GenerateRandStatsForTime(b.N)
	keyMap := GetKeymap()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "benchmark")
	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	b.ResetTimer()

	SaveToCsvWriter(data, keyMap, tmpFile, false)
}

func BenchmarkCsvSaving(b *testing.B) {
	data := GenerateRandStatsForTime(b.N)
	keyMap := GetKeymap()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "benchmark")
	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	b.ResetTimer()

	SaveToCsvWriter(data, keyMap, tmpFile, true)
}

func BenchmarkJSONSavingOnlySum(b *testing.B) {
	data := GenerateRandStatsForTime(b.N)
	keyMap := GetKeymap()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "benchmark")
	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	b.ResetTimer()

	SaveToJSONWriter(data, keyMap, tmpFile, false)
}

func BenchmarkJSONSaving(b *testing.B) {
	data := GenerateRandStatsForTime(b.N)
	keyMap := GetKeymap()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "benchmark")
	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	b.ResetTimer()

	SaveToJSONWriter(data, keyMap, tmpFile, true)
}

func BenchmarkCsvGzSaving(b *testing.B) {
	data := GenerateRandStatsForTime(b.N)
	keyMap := GetKeymap()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "benchmark")
	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	b.ResetTimer()

	SaveToCsvGzWriter(data, keyMap, tmpFile, true)
}

func BenchmarkJSONGzSaving(b *testing.B) {
	data := GenerateRandStatsForTime(b.N)
	keyMap := GetKeymap()

	tmpFile, err := ioutil.TempFile(os.TempDir(), "benchmark")
	if err != nil {
		b.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	b.ResetTimer()

	SaveToJSONGzWriter(data, keyMap, tmpFile, true)
}
