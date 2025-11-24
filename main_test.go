package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	t.Parallel()
	run([]string{"hei", "-d", "-f", ".", "show"})
}
