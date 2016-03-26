// saving_test.go
package main

import (
	"math/rand"
	"testing"
)

func BenchmarkCsvSavingOnlySum(b *testing.B) {
	data := make([]StatForTime, 0)
	keyMap := GetKeymap()
	rnd := rand.New(rand.NewSource(42))
	for i := 0; i < b.N; i++ {
		var curStat StatForTime
		curStat.Init()
		for keyNum := range keyMap {
			if rnd.Float32() > 0.7 {
				curStat.keys[keyNum] = rnd.Intn(5000)
			}
		}
		data = append(data, curStat)
	}
	b.ResetTimer()
	SaveToCsvFile(data, keyMap, "/tmp/bla.csv", true)
}

func BenchmarkCsvSaving(b *testing.B) {
	data := make([]StatForTime, 0)
	keyMap := GetKeymap()
	rnd := rand.New(rand.NewSource(42))
	for i := 0; i < b.N; i++ {
		var curStat StatForTime
		curStat.Init()
		for keyNum := range keyMap {
			if rnd.Float32() > 0.7 {
				curStat.keys[keyNum] = rnd.Intn(5000)
			}
		}
		data = append(data, curStat)
	}
	b.ResetTimer()
	SaveToCsvFile(data, keyMap, "/tmp/bla.csv", false)
}
