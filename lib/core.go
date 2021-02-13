package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// FileExists can be used to check if the given file exists.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// ReadValue reads a value from stdin
func ReadValue(name string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Please specify the %s: ", name)
	scanner.Scan()

	value := strings.TrimSpace(scanner.Text())
	return value
}

func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
