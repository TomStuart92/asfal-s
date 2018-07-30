package log

import "testing"

func TestRecord(t *testing.T) {
	record := Record{Key: "Hello", Value: "World"}
	if record.Key != "Hello" || record.Value != "World" {
		t.Error("Struct Did Not Have Fields Key/Value")
	}
}