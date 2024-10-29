/*
 * @Author: modnarshen
 * @Date: 2024.10.29 11:14:17
 * @Note: Copyrights (c) 2024 modnarshen. All rights reserved.
 */
package json

import (
	"encoding/json"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

var data = map[string]interface{}{
	"name": "John",
	"age":  30,
	"friends": []string{
		"Alice",
		"Bob",
		"Charlie",
	},
}

func BenchmarkStandardJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsoniter(b *testing.B) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkJsoniterFastest(b *testing.B) {
	var json = jsoniter.ConfigFastest
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Error(err)
		}
	}
}
