package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
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

func matchLineHere(line []byte, pattern string) (bool, error) {
	switch {
	case strings.HasPrefix(pattern, `$`) && len(pattern) == 1:
		count := utf8.RuneCount(line)
		return count == 0, nil
	case len(pattern) == 0:
		return true, nil
	case pattern == `$`:
		return len(line) == 0, nil
	case len(pattern) > 1 && pattern[1] == '?':
		return matchZeroOrOne(pattern[0], line, pattern)
	case len(line) == 0:
		return false, nil
	case len(pattern) > 1 && pattern[1] == '+':
		return matchOneOrMore(pattern[0], line, pattern)

	case strings.HasPrefix(pattern, "\\w"):
		char, size := utf8.DecodeRune(line)
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false, nil
		}
		return matchLineHere(line[size:], pattern[2:])

	case strings.HasPrefix(pattern, "."):
		_, size := utf8.DecodeRune(line)
		return matchLineHere(line[size:], pattern[1:])

	case strings.HasPrefix(pattern, "\\d"):
		char, size := utf8.DecodeRune(line)
		if !unicode.IsDigit((char)) {
			return false, nil
		}
		return matchLineHere(line[size:], pattern[2:])

	case strings.HasPrefix(pattern, "("):
		end := strings.IndexByte(pattern, ')')
		middle := strings.IndexByte(pattern, '|')
		first := pattern[1:middle]
		second := pattern[middle+1 : end]
		result, _ := matchLineHere(line, first+pattern[end+1:])
		if result {
			return true, nil
		}
		return matchLineHere(line, second+pattern[end+1:])

	case strings.HasPrefix(pattern, "[^"):
		end := strings.IndexByte(pattern, ']')
		charset := pattern[2:end]
		char, size := utf8.DecodeRune(line)
		if strings.ContainsRune(charset, char) {
			return false, nil
		}
		return matchLineHere(line[size:], pattern[end+1:])

	case strings.HasPrefix(pattern, "["):
		end := strings.IndexByte(pattern, ']')
		charset := pattern[1:end]
		char, size := utf8.DecodeRune(line)
		if !strings.ContainsRune(charset, char) {
			return false, nil
		}
		return matchLineHere(line[size:], pattern[end+1:])
	}

	patternChar, patternCharSize := utf8.DecodeRuneInString(pattern)
	if patternChar == utf8.RuneError {
		return false, fmt.Errorf("bad pattern")
	}

	char, size := utf8.DecodeRune(line)
	if char != patternChar {
		return false, nil
	}
	return matchLineHere(line[size:], pattern[patternCharSize:])
}

func matchLine(line []byte, pattern string) (bool, error) {
	if pattern == "" {
		return true, nil
	}
	if len(line) == 0 {
		return len(pattern) == 0, nil
	}
	if strings.HasPrefix(pattern, "^") {
		return matchLineHere(line, pattern[1:])
	}
	for i := range string(line) {
		ok, err := matchLineHere(line[i:], pattern)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func matchOneOrMore(char byte, line []byte, pattern string) (bool, error) {
	for i := 0; i < len(line); i++ {
		if char != line[i] {
			break
		}
		result, err := matchLine(line[i+1:], pattern[2:])
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil
}

func matchZeroOrOne(char byte, line []byte, pattern string) (bool, error) {
	res, _ := matchLineHere(line, pattern[2:])
	if res {
		return true, nil
	}
	return matchLineHere(line, string(pattern[0])+pattern[2:])
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

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}
