package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/auth"
	"ipmanlk/ani2mal/entries"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	token, err := auth.GetMalAcessCode()

	if err != nil {
		panic(err)
	}

	lists, err := entries.GetMalEntries(token)

	jsonData, _ := json.MarshalIndent(lists, "", " ")

	os.WriteFile("lala.json", jsonData, 0644)
}
