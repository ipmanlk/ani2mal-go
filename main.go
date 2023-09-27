package main

import (
	"fmt"
	"ipmanlk/ani2mal/anilist"
)

func main() {
	fmt.Println("Hello, World!")

	anilist.PerformAuth()

	// jsonData, _g := json.MarshalIndent(res, "", " ")

	// os.WriteFile("lala.json", jsonData, 0644)
}
