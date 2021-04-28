package main

import "testing"

func TestSayHello(t *testing.T) {
	output := SayHello()
	expected := "Hello world!"

	if output != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output)
	}
}
