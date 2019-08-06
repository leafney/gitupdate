package main

import (
	"strings"
	"testing"
)

func TestA(t *testing.T) {
	baBranch := "origin/release-2.9.0"
	// newBranch := strings.TrimLeft(baBranch, "origin//")
	newBranch := strings.TrimPrefix(baBranch, "origin/")
	t.Log(newBranch)
}
