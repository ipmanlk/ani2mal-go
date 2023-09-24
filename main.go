package main

import (
	"encoding/json"
	"fmt"
	"ipmanlk/ani2mal/entries"
	"os"
)

func main() {
	fmt.Println("Hello, World!")

	res, err := entries.GetList("CrystalBullet", "ANIME")

	if err != nil {
		panic(err)
	}

	jsonData, err := json.MarshalIndent(res, "", " ")

	os.WriteFile("lala.json", jsonData, 0644)
}
