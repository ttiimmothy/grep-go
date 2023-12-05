package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func IsDigit(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func matchLine(line string, pattern string, matchStart bool) (bool, error) {
	if len(pattern) == 0 {
		return true, nil
	}
	if len(line) == 0 {
		return len(pattern) == 0, nil
	}
	if strings.HasPrefix(pattern, "^") {
		return matchLine(line, pattern[1:], true)
	}
	if strings.HasPrefix(pattern, "[") {
		end := strings.IndexByte(pattern, ']')
		charset := pattern[1:end]
		if charset[0] == '^' {
			for _, c := range charset {
				if strings.HasPrefix(line, string(c)) {
					return false, nil
				}
			}
			return matchLine(line[1:], pattern[end+1:], matchStart)
		} else {
			for _, c := range charset {
				if strings.HasPrefix(line, string(c)) {
					result, _ := matchLine(line[1:], pattern[end+1:], true)
					if result {
						return true, nil
					}
				}
			}
		}
		return false, nil
	}

	char := string(line[0])
	if strings.HasPrefix(pattern, "\\w") {
		if IsLetter(char) || IsDigit(char) {
			return matchLine(line[1:], pattern[2:], matchStart)
		} else {
			if matchStart {
				return false, nil
			} else {
				return matchLine(line[1:], pattern, matchStart)
			}
		}
	} else if strings.HasPrefix(pattern, "\\d") {
		if IsDigit(char) {
			return matchLine(line[1:], pattern[2:], matchStart)
		} else {
			if matchStart {
				return false, nil
			} else {
				return matchLine(line[1:], pattern, matchStart)
			}
		}
	} else {
		if string(pattern[0]) == char {
			return matchLine(line[1:], pattern[1:], matchStart)
		} else {
			if matchStart {
				return false, nil
			} else {
				return matchLine(line[1:], pattern, matchStart)
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(string(line), pattern, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}
