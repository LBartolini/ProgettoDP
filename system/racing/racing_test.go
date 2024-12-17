package main

import (
	"log"
	"testing"
)

func TestExample(t *testing.T) {
	log.Printf("This is a test")
	t.Fail()
}
