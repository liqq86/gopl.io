package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type WordsCounter int

func (wc *WordsCounter) Write(data []byte) (n int, err error) {
	s := string(data)
	scanner := bufio.NewScanner(strings.NewReader(s))
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	n = count
	return
}

func main() {
	var wc WordsCounter
	n, err := wc.Write([]byte("hello, this is liqq speaking"))
	fmt.Println(n, err)
}
