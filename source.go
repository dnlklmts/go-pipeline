package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetDataStream() (<-chan int, <-chan bool) {
	done := make(chan bool)
	dataStream := make(chan int)

	go func() {
		defer close(done)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter an integer (or type \"exit\").")

		for {
			fmt.Print("Input: ")
			scanner.Scan()

			input := scanner.Text()
			if strings.EqualFold(input, "exit") {
				fmt.Println("Completed.")
				return
			}

			data, err := strconv.Atoi(strings.TrimSpace(input))
			if err != nil {
				fmt.Println("Skipped. Only integers are allowed.")
				continue
			}

			dataStream <- data
		}
	}()

	return dataStream, done
}
