package main

import (
	"testing"
)

func BenchmarkProcessFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProcessFiles("tmp", HelperUtils{loadDataToStruct, GetUniqueKey, hardCheck, getAllFiles})
	}
}

func BenchmarkProcessFilesUsingChannels(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ProcessFilesWithoutMutex("tmp", HelperUtils{loadDataToStruct, GetUniqueKey, hardCheck, getAllFiles})
	}
}
