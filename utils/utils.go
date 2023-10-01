package utils

import (
	"bufio"
	"os"
)

func GetStrInput() string {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	input, err := reader.ReadString('\n')

	// 	if err != nil {
	// 		log.Fatal(err)
	// }

	// 	return strings.TrimSpace(input)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	return input
}
