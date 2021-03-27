package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadInputStreamBuffer(reader *bufio.Reader) ([]string, error) {
	var fileArgs []string
	for {
		input, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		line := string(input)
		if line != "" {
			if strings.Contains(line, string(os.PathSeparator)) && fileExists(line) {
				fileArgs = append(fileArgs, line)
			} else {
				return nil, newUnKnownInputStreamError(line)
			}
		}
	}

	return fileArgs, nil
}
