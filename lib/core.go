package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
