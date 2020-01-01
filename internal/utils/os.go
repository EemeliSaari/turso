package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func FileLines(path string) []string {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)

	lines := []string{}
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		} else {
			line = strings.Trim(line, " \n\t\r")
			lines = append(lines, line)
		}
	}
	fmt.Println(lines)
	return lines
}
