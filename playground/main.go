package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	// This is a placeholder for the main function.
	// You can add your code here to execute when the program runs.
	filename := "./test/test.log"
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	_ = reader
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			println(line)
		}
		if err == io.EOF {
			break
		}
	}
}
