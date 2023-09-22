package utils

import (
	"bufio"
	"log"
	"os"
)

func GetStrInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	return input
}
