// saving_test.go
package main

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func GenerateRandStatsForTime(N int) []StatForTime {
	data := make([]StatForTime, 0)
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

	SaveToCsvWriter(data, keyMap, tmpFile, true)
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

	SaveToCsvWriter(data, keyMap, tmpFile, false)
}
