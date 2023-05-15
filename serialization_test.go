package PinkRPC

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
)

type Person struct {
	Id     string
	Name   string
	Age    int
	Parent []Person
}

var (
	PersonA = Person{
		Id:     "lowe234*7234#los7823",
		Name:   "BaiYuLong",
		Age:    22,
		Parent: nil,
	}
	PersonB = Person{
		Id:     "lowe234*7234#los7823",
		Name:   "LiuHongJiang",
		Age:    22,
		Parent: []Person{PersonA},
	}
)

func TestPob(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(PersonB)
	fmt.Printf("gob length %d\n", buf.Len())
	var person Person
	dec := gob.NewDecoder(&buf)
	_ = dec.Decode(&person)
	fmt.Printf("gob result %v\n", person)
}

func TestJson(t *testing.T) {
	js, _ := json.Marshal(PersonB)
	fmt.Printf("json length %d\n", len(js))
	var person Person
	_ = json.Unmarshal(js, &person)
	fmt.Printf("json result %v\n", person)
}

func BenchmarkGobEncoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		_ = enc.Encode(PersonB)
	}
}

func BenchmarkGobDecoder(b *testing.B) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(PersonB)
	var person Person
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec := gob.NewDecoder(&buf)
		_ = dec.Decode(&person)
	}
}

func BenchmarkJsonEncoder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(PersonB)
	}
}

func BenchmarkJsonDecoder(b *testing.B) {
	js, _ := json.Marshal(PersonB)
	var person Person
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(js, &person)
	}
}
