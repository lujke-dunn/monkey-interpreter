package object

import "testing"

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is luke"}
	diff2 := &String{Value: "My name is luke"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hashkeys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have the same hashkeys")
	}
}

